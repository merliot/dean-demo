package metro

import (
	"embed"
	"net"
	"net/http"

	"github.com/merliot/dean"
)

//go:embed css js images index.html
var fs embed.FS

type Metro struct {
	dean.Thing
	dean.ThingMsg
	CPUFreq float64
	Mac     string
	Ip      net.IP
	Input   bool
	runChan chan *dean.Msg
}

func New(id, model, name string) dean.Thinger {
	println("NEW METRO")
	return &Metro{
		Thing:   dean.NewThing(id, model, name),
		runChan: make(chan *dean.Msg, 10),
	}
}

func (m *Metro) saveState(msg *dean.Msg) {
	msg.Unmarshal(m)
}

func (m *Metro) getState(msg *dean.Msg) {
	m.Path = "state"
	msg.Marshal(m).Reply()
}

func (m *Metro) broadcast(msg *dean.Msg) {
	msg.Unmarshal(m).Broadcast()
}

func (m *Metro) run(msg *dean.Msg) {
	msg.Unmarshal(m).Broadcast()
	if m.IsMetal() {
		m.runChan <- msg
	}
}

func (m *Metro) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"state":     m.saveState,
		"get/state": m.getState,
		"tx":        m.run,
		"input":     m.broadcast,
		"reset":     m.run,
	}
}

func (m *Metro) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.ServeFS(fs, w, r)
}
