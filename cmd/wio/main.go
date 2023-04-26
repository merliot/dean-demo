package main

import (
	"github.com/merliot/dean"
	"github.com/merliot/dean-demo/wio"
)

func main() {
	thing := wio.New("demo-wio-01", "demo-wio", "wio")
	server := dean.NewServer(thing)
	server.DialWebSocket("", "", "wss://demo.merliot.net/ws/1500", thing.Announce())
	server.Run()
}
