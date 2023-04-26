//go:build tinygo

package metro

import (
	"machine"
	"time"

	"github.com/merliot/dean"
	"github.com/merliot/dean/tinynet"
	"github.com/merliot/dean/lora/lorae5"
)

type runMsg struct {
	Path string
}

type txMsg struct {
	Path string
	Tx   string
}

type inputMsg struct {
	Path  string
	Input bool
}

func (m *Metro) Run(i *dean.Injector) {
	var msg dean.Msg

	m.CPUFreq = float64(machine.CPUFrequency()) / 1000000.0
	mac, _ := tinynet.GetHardwareAddr()
	m.Mac = mac.String()
	m.Ip, _ = tinynet.GetIPAddr()

	// Input on GPIO D2
	d2 := machine.D2
	d2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	m.Input = d2.Get()

	m.Path = "update"
	i.Inject(msg.Marshal(m))

	lora := lorae5.New(machine.UART1, machine.UART_TX_PIN, machine.UART_RX_PIN, 9600)
	lora.Init()

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			input := !d2.Get()
			if input != m.Input {
				var imsg = inputMsg{Path: "input", Input: input}
				i.Inject(msg.Marshal(&imsg))
			}
		case msg := <-m.runChan:
			var rmsg runMsg
			msg.Unmarshal(&rmsg)
			switch rmsg.Path {
			case "tx":
				var tmsg txMsg
				msg.Unmarshal(&tmsg)
				err := lora.Tx([]byte(tmsg.Tx), 1000)
				if err != nil {
					println(err.Error())
				}
			case "reset":
				machine.CPUReset()
			}
		}
	}
}
