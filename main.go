// Source code file, created by Developer@YANYINGSONG.

package main

import (
	"arena"
	"net/http"
	"r/pkg/connect"
	"r2/pkg/generic2/chars/cat"
	"r2/pkg/generic2/signalgroup"
	"strings"
	"time"
)

var _ = arena.NewArena()

func main() {
	// rpcpd.DebugOn()

	for _, line := range strings.Split(cat.String("ua.txt"), "\n") {
		if line != "" {
			connect.AccessUA[line] = nil
		}
	}
	connect.InitDeviceAlias()

	signalgroup.Async(connect.RunServer)
	signalgroup.Async(webserver)
	signalgroup.Async(func() error {
		time.Sleep(3 * time.Hour)
		return nil
	})

	signalgroup.Wait(nil)
}

func webserver() error {
	http.HandleFunc("/call", connect.CallHandler)
	http.HandleFunc("/ls", connect.DeviceList)
	err := http.ListenAndServe(":90", nil)
	if err != nil {
		panic(err)
	}

	return nil
}
