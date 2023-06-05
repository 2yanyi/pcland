// Source code file, created by Developer@YANYINGSONG.

package main

import (
	"arena"
	"fmt"
	"library/generic/errcause"
	"os"
	methods "r/client/methods"
	"r/server"
	"rpcpd"
	"rpcpd/seqx"
	"strings"
	"time"
)

var _ = arena.NewArena()
var client *rpcpd.Power

func runClient(serverAddr string) func() {
	seqx.Init(1)

	var err error
	client, err = rpcpd.Connect(serverAddr + server.ListenPort)
	if err != nil {
		fmt.Printf("failed: %s\n", err.Error())
	}

	return func() {
		if client != nil {
			if client.Client != nil {
				client.Client.Channel.Close()
			}
		}
	}
}

func delivery(serverAddr string) {
	defer runClient(serverAddr)()
	if client == nil {
		return
	}
	rpcpd.AddFunction(methods.SignPing, methods.Ping)
	rpcpd.AddFunction(methods.SignGetDeviceInfo, methods.GetDeviceInfo)
	rpcpd.AddFunction(methods.SignRestart, methods.Restart)

	// 保持长连接
	connectionLive()
}

func connectionLive() {
	defer errcause.Recover()
	if err := client.ConnectionProcessor(nil); err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func main() {
	// rpcpd.DebugOn()

	imageArgs := strings.Split(os.Args[0], "-@")
	serverAddr := imageArgs[len(imageArgs)-1]
	serverAddr = serverAddr[:strings.LastIndexByte(serverAddr, '.')]
	fmt.Printf("serverAddr=%s\n", serverAddr)

	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	delivery(serverAddr)
	for range t.C {
		delivery(serverAddr)
	}
}
