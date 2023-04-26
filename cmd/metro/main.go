package main

import (
	"github.com/merliot/dean"
	"github.com/merliot/dean/demo/metro"
)

func main() {
	thing := metro.New("demo-metro-01", "demo-metro", "metro")
	server := dean.NewServer(thing)
	server.DialWebSocket("", "", "wss://demo.merliot.net/ws/1500", thing.Announce())
	server.Run()
}
