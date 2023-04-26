//go:build tinygo

package matrix

import (
	"machine"

	"github.com/merliot/dean"
	"github.com/merliot/dean/tinynet"
	"github.com/merliot/dean-lib/lora/lorae5"
)

type rxMsg struct {
	Path   string
	Rx     string
}

type relayMsg struct {
	Path  string
	Relay bool
}

func (m *Matrix) tinyGoIsGreat() bool {
	if m.Rx == "TinyGo" {
		m.lastRx = ""
	}
	m.lastRx += m.Rx
	relay := (m.lastRx == "TinyGoIsGreat!")
	if relay != m.Relay {
		m.Relay = relay
		return true
	}
	return false
}

func (m *Matrix) Run(i *dean.Injector) {
	var msg dean.Msg
	var loraOut = make(chan []byte)

	m.CPUFreq = float64(machine.CPUFrequency()) / 1000000.0
	mac, _ := tinynet.GetHardwareAddr()
	m.Mac = mac.String()
	m.Ip, _ = tinynet.GetIPAddr()

	m.Path = "update"
	i.Inject(msg.Marshal(m))

	relay := machine.A3
	relay.Configure(machine.PinConfig{Mode: machine.PinOutput})

	lora := lorae5.New(machine.UART1, machine.UART1_TX_PIN, machine.UART1_RX_PIN, 9600)
	lora.Init()
	go lora.RxPoll(loraOut, 2000)

	for {
		select {
		case pkt := <-loraOut:
			m.Rx = string(pkt)
			rmsg := rxMsg{Path: "rx", Rx: m.Rx}
			i.Inject(msg.Marshal(&rmsg))
			if m.tinyGoIsGreat() {
				relay.Set(m.Relay)
				rmsg := relayMsg{Path: "relay", Relay: m.Relay}
				i.Inject(msg.Marshal(&rmsg))
			}
		case <-m.runChan:
			machine.CPUReset()
		}
	}
}
