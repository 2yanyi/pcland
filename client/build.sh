#!/bin/bash

# shellcheck disable=SC1035
if !(type go &> /dev/null); then
  echo 'go: command not found'
  exit
fi

# Default config
export GOEXPERIMENT=arenas
export GOARCH=amd64

# GNU/Linux
echo "Build GNU/Linux"
ldflags='-linkmode "external" -extldflags "-static"'
export GOOS=linux
GOAMD64=v3 \
GOARCH=amd64   go build -ldflags "${ldflags}" -o ../bin/linux.amd64
GOARCH=amd64   go build -ldflags "${ldflags}" -o ../bin/linux.amd64v1
GOARCH=386     go build                       -o ../bin/linux.386
GOARCH=arm64   go build                       -o ../bin/linux.arm64
GOARCH=loong64 go build                       -o ../bin/linux.loong64
GOARCH=riscv64 go build                       -o ../bin/linux.riscv64

# Windows
echo "Build Windows"
export GOOS=windows
GOAMD64=v3 \
GOARCH=amd64   go build -ldflags "-H windowsgui" -o ../bin/windows.amd64
GOARCH=amd64   go build -ldflags "-H windowsgui" -o ../bin/windows.amd64v1
GOARCH=386     go build -ldflags "-H windowsgui" -o ../bin/windows.386
GOARCH=arm64   go build -ldflags "-H windowsgui" -o ../bin/windows.arm64

# macOS
echo "Build macOS"
export GOOS=darwin
GOARCH=amd64 go build -o ../bin/apple.amd64
GOARCH=arm64 go build -o ../bin/apple.arm64
