#!/bin/bash

# shellcheck disable=SC1035
if !(type go &> /dev/null); then
  echo 'go: command not found'
  exit
fi

# Default config
export GOEXPERIMENT=arenas
export GOARCH=amd64

# Exe name
EXE='PCland-@211.149.130.119'

# GNU/Linux
echo "Build GNU/Linux"
ldflags='-linkmode "external" -extldflags "-static"'
export GOOS=linux
GOAMD64=v3 \
GOARCH=amd64   go build -ldflags "${ldflags}" -o ../bin/linux/${EXE}.amd64
GOARCH=amd64   go build -ldflags "${ldflags}" -o ../bin/linux/${EXE}.amd64v1
GOARCH=386     go build                       -o ../bin/linux/${EXE}.386
GOARCH=arm64   go build                       -o ../bin/linux/${EXE}.arm64
GOARCH=loong64 go build                       -o ../bin/linux/${EXE}.loong64
GOARCH=riscv64 go build                       -o ../bin/linux/${EXE}.riscv64

# Windows
echo "Build Windows"
export GOOS=windows
GOAMD64=v3 \
GOARCH=amd64   go build -ldflags "-H windowsgui" -o ../bin/windows/${EXE}.amd64.exe
GOARCH=amd64   go build -ldflags "-H windowsgui" -o ../bin/windows/${EXE}.amd64v1.exe
GOARCH=386     go build -ldflags "-H windowsgui" -o ../bin/windows/${EXE}.386.exe
GOARCH=arm64   go build -ldflags "-H windowsgui" -o ../bin/windows/${EXE}.arm64.exe

# macOS
echo "Build macOS"
export GOOS=darwin
GOARCH=amd64 go build -o ../bin/apple/${EXE}.amd64
GOARCH=arm64 go build -o ../bin/apple/${EXE}.arm64
