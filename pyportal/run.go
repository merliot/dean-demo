//go:build !tinygo

package pyportal

import "github.com/merliot/dean"

func (p *Pyportal) Run(i *dean.Injector) {
	for {
		select {
		case <-p.runChan:
		}
	}
}
