#!/bin/bash

# shellcheck disable=SC1035
if !(type go &> /dev/null); then
    echo 'go: command not found'
    exit
fi

# go flags
export GOEXPERIMENT=arenas
export GOARCH=amd64
export GOOS=linux
export GOAMD64=v2

go build -o bin/pcland-svr
