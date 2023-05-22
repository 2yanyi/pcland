// Source code file, created by Developer@YANYINGSONG.

package main

import (
	"arena"
	"fmt"
	"os"
	"r/pkg/connect"
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
	client, err = rpcpd.Connect(serverAddr + connect.ServerListenPort)
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
	rpcpd.AddFunction(connect.FunctionSignGetDeviceInfo, connect.GetDeviceInfo)
	rpcpd.AddFunction(connect.FunctionSignShutdown, connect.Shutdown)
	rpcpd.AddFunction(connect.FunctionSignRestart, connect.Restart)

	// 保持长连接
	client.ConnectionProcessor(nil)
}

func main() {
	// rpcpd.DebugOn()

	name := strings.Split(os.Args[0], "_")
	serverAddr := name[len(name)-1]
	if strings.HasSuffix(serverAddr, ".exe") {
		serverAddr = strings.TrimSuffix(serverAddr, ".exe")
	}
	if strings.HasSuffix(serverAddr, ".m1") {
		serverAddr = strings.TrimSuffix(serverAddr, ".m1")
	}
	fmt.Printf("serverAddr=%s\n", serverAddr)

	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	delivery(serverAddr)
	for range t.C {
		delivery(serverAddr)
	}
}
