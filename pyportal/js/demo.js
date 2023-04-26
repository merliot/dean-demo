var state
var conn
var online = false

function neo(letter, value) {
	state.NeoColor[letter] = parseInt(value)
	conn.send(JSON.stringify({Path: "neo", NeoColor: state.NeoColor}))
	showNeoPixel()
}

function showLight() {
	let light = document.getElementById("light")
	light.value = ""
	percent = state.Light * 100 / 65535
	light.value = "Light Intensity: " + percent.toFixed(3)
}

function showTemp() {
	let tempc = document.getElementById("tempc")
	tempc.value = ""
	tempc.value = "Temperature:     " + state.TempC.toFixed(2) + "(C)"
}

function showSystem() {
	let system = document.getElementById("system")
	system.value = ""
	system.value += "CPU Frequency:   " + state.CPUFreq + "Mhz\r\n"
	system.value += "MAC Address:     " + state.Mac + "\r\n"
	system.value += "IP Address:      " + state.Ip
}

function showNeoPixel() {
	let circle = document.getElementById("neopixel");
	let r = state.NeoColor["R"]
	let g = state.NeoColor["G"]
	let b = state.NeoColor["B"]
	let a = state.NeoColor["A"] / 255.0
	let colorString = "rgba(" + r + "," + g + "," + b + "," + a + ")";
	circle.setAttribute("fill", colorString);
}

function showNeo() {
	let neoR = document.getElementById("neoR")
	let neoG = document.getElementById("neoG")
	let neoB = document.getElementById("neoB")
	neoR.value = state.NeoColor["R"]
	neoG.value = state.NeoColor["G"]
	neoB.value = state.NeoColor["B"]
	showNeoPixel()
}

function reset() {
	state.NeoColor["R"] = 0
	state.NeoColor["G"] = 0
	state.NeoColor["B"] = 0
	showNeo()
	conn.send(JSON.stringify({Path: "neo", NeoColor: state.NeoColor}))
	conn.send(JSON.stringify({Path: "reset"}))
}

function show() {
	overlay = document.getElementById("overlay")
	overlay.style.display = online ? "none" : "block"
	showSystem()
	showLight()
	showTemp()
	showNeo()
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
		console.log('pyportal', msg)

		switch(msg.Path) {
		case "state":
			online = true
			// fall-thru
		case "update":
			state = msg
			show()
			break
		case "light":
			state.Light = msg.Light
			showLight()
			break
		case "tempc":
			state.TempC = msg.TempC
			showTemp()
			break
		case "neo":
			state.NeoColor = msg.NeoColor
			showNeo()
			break
		}
	}
}

