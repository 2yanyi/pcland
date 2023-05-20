#!/bin/bash

# shellcheck disable=SC1035
if !(type go &> /dev/null); then
  echo 'go: command not found'
  exit
fi

# Default config
export GOEXPERIMENT=arenas
export GOARCH=amd64
export GOAMD64=v3

# Exe name
APP='pcland'
IP='211.149.130.119'

# GNU/Linux
export GOOS=linux
           go build -ldflags '-linkmode "external" -extldflags "-static"' -o ../bin/${APP}64_${IP}
GOAMD64=v1 go build -ldflags '-linkmode "external" -extldflags "-static"' -o ../bin/${APP}64o_${IP}

# Windows
export GOOS=windows
           go build -ldflags "-H windowsgui" -o ../bin/${APP}64_${IP}.exe
GOARCH=386 go build -ldflags "-H windowsgui" -o ../bin/${APP}32_${IP}.exe

# macOS
export GOOS=darwin
GOARCH=arm64 go build -o ../bin/${APP}64_${IP}.m1
