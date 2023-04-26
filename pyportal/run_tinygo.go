//go:build tinygo

package pyportal

import (
	"image/color"
	"machine"
	"time"

	"github.com/merliot/dean"
	"github.com/merliot/dean/tinynet"
	"tinygo.org/x/drivers/adt7410"
	"tinygo.org/x/drivers/ws2812"
)

type lightMsg struct {
	Path string
	Light uint16
}

type tempMsg struct {
	Path  string
	TempC float32
}

type runMsg struct {
	Path string
}

func (p *Pyportal) Run(i *dean.Injector) {
	var msg dean.Msg

	p.CPUFreq = float64(machine.CPUFrequency()) / 1000000.0
	mac, _ := tinynet.GetHardwareAddr()
	p.Mac = mac.String()
	p.Ip, _ = tinynet.GetIPAddr()

	// Ambient light sensor
	lightSensor := machine.ADC{machine.A2}
	lightSensor.Configure(machine.ADCConfig{})
	p.Light = lightSensor.Get()

	// ADT7410 temperature sensor
	machine.I2C0.Configure(machine.I2CConfig{})
	tempSensor := adt7410.New(machine.I2C0)
	tempSensor.Configure()
	p.TempC = tempSensor.ReadTempC()

	// Neo Pixel
	neo := machine.WS2812
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(neo)
	ws.WriteColors([]color.RGBA{p.NeoColor})

	p.Path = "update"
	i.Inject(msg.Marshal(p))

	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			light := lightSensor.Get()
			if light != p.Light {
				p.Light = light
				lmsg := lightMsg{Path: "light", Light: light}
				i.Inject(msg.Marshal(&lmsg))
			}
			tempC := tempSensor.ReadTempC()
			if tempC != p.TempC {
				p.TempC = tempC
				tmsg := tempMsg{Path: "tempc", TempC: tempC}
				i.Inject(msg.Marshal(&tmsg))
			}
		case msg := <-p.runChan:
			var rmsg runMsg
			msg.Unmarshal(&rmsg)
			switch rmsg.Path {
			case "neo":
				// Alpha channel is not supported by WS2812
				ws.WriteColors([]color.RGBA{p.NeoColor})
			case "reset":
				machine.CPUReset()
			}
		}
	}
}
