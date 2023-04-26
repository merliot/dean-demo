# Dean Demo

## Install

```
go install -C cmd/demo
```

Give perms to serve on port :443

```
sudo setcap CAP_NET_BIND_SERVICE=+eip ~/go/bin/demo
```

## Run

Run demo as https web server on \<host\> on port :443.

```
~/go/bin/demo -host <host>
```

## Building Firmware

Build the firmware for each target using TinyGo\*:

```
tinygo flash -monitor -target pyportal -stack-size 4KB cmd/pyportal/main.go
tinygo flash -monitor -target nano-rp2040 -stack-size 4kB cmd/connect/main.go
tinygo flash -monitor -target metro-m4-airlift -stack-size 4KB cmd/metro/main.go
tinygo flash -monitor -target matrixportal-m4 -stack-size 4KB cmd/matrix/main.go
```

\* Requires netdev patch to tinygo.
