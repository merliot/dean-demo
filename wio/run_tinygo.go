//go:build tinygo

package wio

import (
	"machine"
	"time"

	"github.com/merliot/dean"
	"github.com/merliot/dean/tinynet"
)

func (w *Wio) Run(i *dean.Injector) {
	var msg dean.Msg

	ticker := time.NewTicker(time.Second)

	w.CPUFreq = float64(machine.CPUFrequency()) / 1000000.0
	mac, _ := tinynet.GetHardwareAddr()
	w.Mac = mac.String()
	w.Ip, _ = tinynet.GetIPAddr()

	w.Path = "update"
	i.Inject(msg.Marshal(w))

	for {
		select {}
	}
}
