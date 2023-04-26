//go:build tinygo

package connect

import (
	"crypto/rand"
	"machine"
	"time"

	"github.com/merliot/dean"
	"github.com/merliot/dean/tinynet"
	"tinygo.org/x/drivers/bh1750"
)

type luxMsg struct {
	Path string
	Lux  int32
}

type tempMsg struct {
	Path  string
	TempC int32
}

func (c *Connect) Run(i *dean.Injector) {
	var msg dean.Msg

	c.CPUFreq = float64(machine.CPUFrequency()) / 1000000.0
	mac, _ := tinynet.GetHardwareAddr()
	c.Mac = mac.String()
	c.Ip, _ = tinynet.GetIPAddr()
	c.TempC = machine.ReadTemperature() / 1000

	// Relay on GPIO D2
	relay := machine.D2
	relay.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// BH1750 Light Intensity
	machine.I2C0.Configure(machine.I2CConfig{})
	sensor := bh1750.New(machine.I2C0)
	sensor.Configure()
	c.Lux = sensor.Illuminance()

	setRelay := func() {
		if 650000 <= c.Lux && c.Lux <= 700000 {
			relay.High()
		} else {
			relay.Low()
		}
	}
	setRelay()

	c.Path = "update"
	i.Inject(msg.Marshal(c))

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			temp := machine.ReadTemperature() / 1000
			if temp != c.TempC {
				c.TempC = temp
				tmsg := tempMsg{Path: "tempc", TempC: temp}
				i.Inject(msg.Marshal(&tmsg))
			}
			lux := sensor.Illuminance()
			if lux != c.Lux {
				c.Lux = lux
				setRelay()
				lmsg := luxMsg{Path: "lux", Lux: lux}
				i.Inject(msg.Marshal(&lmsg))
			}
		case <-c.runChan:
			machine.CPUReset()
		}
	}
}

// TODO: remove below when RNG is working on rp2040

func init() {
	rand.Reader = &reader{}
}

type reader struct{}

func (r *reader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return
	}
	var randomByte uint32
	for i := range b {
		if i%4 == 0 {
			randomByte, err = machine.GetRNG()
			if err != nil {
				return n, err
			}
		} else {
			randomByte >>= 8
		}
		b[i] = byte(randomByte)
	}
	return len(b), nil
}
