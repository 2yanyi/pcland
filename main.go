// Source code file, created by Developer@YANYINGSONG.

package main

import (
	"arena"
	"net/http"
	"r/pages"
	"r/server"
	"r2/pkg/generic2/chars"
	"r2/pkg/generic2/signalgroup"
	"r2/pkg/utils/console"
	"strings"
	"time"
)

var _ = arena.NewArena()

func main() {
	// rpcpd.DebugOn()

	defer console.New(&console.Options{
		Info: true, Debug: true, Warning: true, Error: true, Print: true,
		LogFileSizeMB: 100,
		MaxBackups:    10,
		Filename:      "log/execution.log",
	}).Wait()

	for _, line := range strings.Split(chars.CatString("ua.txt"), "\n") {
		if line != "" {
			pages.AccessUA[line] = nil
		}
	}
	server.InitDeviceAlias()

	signalgroup.Async(server.RunServer)
	signalgroup.Async(webserver)
	signalgroup.Async(func() error {
		time.Sleep(3 * time.Hour)
		return nil
	})

	signalgroup.Wait(nil)
}

func webserver() error {
	http.HandleFunc("/get", pages.DownloadHandler)
	http.HandleFunc("/call", pages.CallHandler)
	http.HandleFunc("/ls", pages.DeviceList)
	err := http.ListenAndServe(":90", nil)
	if err != nil {
		panic(err)
	}

	return nil
}
