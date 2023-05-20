// Source code file, created by Developer@YANYINGSONG.

package main

import (
	"arena"
	"net/http"
	"r/pkg/connect"
	"r2/pkg/generic2/chars/cat"
	"strings"
)

func main() {
	// rpcpd.DebugOn()

	for _, line := range strings.Split(cat.String("ua.txt"), "\n") {
		if line != "" {
			connect.AccessUA[line] = nil
		}
	}
	connect.InitDeviceAlias()

	go connect.RunServer()
	http.HandleFunc("/call", connect.CallHandler)
	http.HandleFunc("/ls", connect.DeviceList)
	err := http.ListenAndServe(":90", nil)
	if err != nil {
		panic(err)
	}
}

var mem = arena.NewArena()
