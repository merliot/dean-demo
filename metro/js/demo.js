var state
var conn
var online = false

function tx(msg) {
	conn.send(JSON.stringify({Path: "tx", Tx: msg}))
}

function showSystem() {
	let system = document.getElementById("system")
	system.value = ""
	system.value += "CPU Frequency:   " + state.CPUFreq + "Mhz\r\n"
	system.value += "MAC Address:     " + state.Mac + "\r\n"
	system.value += "IP Address:      " + state.Ip
}

function showGopher() {
	gopher = document.getElementById("gopher")
	gopher.style.visibility = (state.Input ? 'visible' : 'hidden')
	gopher.src = "images/gopher-workout.gif"
}

function show() {
	overlay = document.getElementById("overlay")
	overlay.style.display = online ? "none" : "block"
	showSystem()
	showGopher()
}

function reset() {
	state.Input = false
	showGopher()
	conn.send(JSON.stringify({Path: "input", Input: false}))
	conn.send(JSON.stringify({Path: "reset"}))
}

function run(ws) {

	console.log('connecting...')
	conn = new WebSocket(ws)

	conn.onopen = function(evt) {
		console.log("open")
		conn.send(JSON.stringify({Path: "get/state"}))
	}

	conn.onclose = function(evt) {
		console.log("close")
		online = false
		show()
		setTimeout(run(ws), 1000)
	}

	conn.onerror = function(err) {
		console.log("error", err)
		conn.close()
	}

	conn.onmessage = function(evt) {
		msg = JSON.parse(evt.data)
		console.log('metro', msg)

		switch(msg.Path) {
		case "state":
			online = true
			state = msg
			show()
			break
		case "input":
			state.Input = msg.Input
			showGopher()
			break
		}
	}
}

