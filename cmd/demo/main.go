package main

import (
	"log"

	"github.com/merliot/dean"
	"github.com/merliot/dean-demo"
	"github.com/merliot/dean-demo/connect"
	"github.com/merliot/dean-demo/matrix"
	"github.com/merliot/dean-demo/metro"
	"github.com/merliot/dean-demo/pyportal"
)

func main() {
	demo := demo.New("demo01", "demo", "demo1").(*demo.Demo)

	server := dean.NewServer(demo)
	server.MaxSockets(100)

	demo.Register("demo-connect", connect.New)
	demo.Register("demo-matrix", matrix.New)
	demo.Register("demo-metro", metro.New)
	demo.Register("demo-pyportal", pyportal.New)

	log.Fatal(server.ServeTLS("demo.merliot.net"))
}
