package demo

import (
	"embed"
	"net/http"

	"github.com/merliot/dean"
	"github.com/merliot/dean-lib/hub"
)

//go:embed css js index.html
var fs embed.FS

type Demo struct {
	*hub.Hub
}

func New(id, model, name string) dean.Thinger {
	println("NEW DEMO")
	return &Demo{
		Hub: hub.New(id, model, name).(*hub.Hub),
	}
}

func (d *Demo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.ServeFS(fs, w, r)
}
