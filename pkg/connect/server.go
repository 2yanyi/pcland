// Source code file, created by Developer@YANYINGSONG.

package connect

import (
	"encoding/json"
	"fmt"
	"os"
	"r2/pkg/generic2/chars/cat"
	"r2/pkg/generic2/errcause"
	"rpcpd"
	"rpcpd/seqx"
	"strings"
	"time"
)

const ServerListenPort = ":11299"

type DeviceInfo struct {
	Power  int    `json:"power,omitempty"`
	Target string `json:"target,omitempty"`
	OS     string `json:"os,omitempty"`
	Model  string `json:"model,omitempty"`
	Addr   string `json:"addr,omitempty"`
	Tid    string `json:"tid,omitempty"`

	Conn *rpcpd.Conn `json:"-"`
}

var ConnectionPool = make(map[string]*DeviceInfo, 10)

var Power *rpcpd.Power

var deviceAlias = make(map[string]string)

func InitDeviceAlias() {
	data := cat.Bytes("alias.json")
	if err := json.Unmarshal(data, &deviceAlias); err != nil {
		println(err.Error())
	}
}

func RunServer() error {
	seqx.Init(2)

	var task func() error
	var powerErr error

	Power, task, powerErr = rpcpd.Listen(ServerListenPort, "tls/cert.pem", "tls/key.pem")
	if powerErr != nil {
		panic(powerErr)
	}
	defer Power.Listener.Close()

	__task := func() {
		if err := task(); err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(-1)
		}
	}
	_ = __task
	// go __task()

	for {
		conn, err := Power.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		go func() {
			defer errcause.Recover()
			Power.ConnectionProcessor(conn)
			deleteConnection(conn)
		}()

		go func() {
			if !addConnection(conn) {
				conn.Channel.Close()
			}
		}()
	}
}

func deleteConnection(conn *rpcpd.Conn) {
	delete(ConnectionPool, conn.ID)
}

func addConnection(conn *rpcpd.Conn) bool {
	fmt.Printf("connected client %s\n", conn.Channel.RemoteAddr().String())
	time.Sleep(3 * time.Second)

	// call client get device info.
	resp, err := Power.Call(conn, []byte(FunctionSignGetDeviceInfo), nil)
	if err != nil {
		fmt.Printf("failed call: %s\n", err)
		return false
	}
	if resp == nil {
		// if head.Code != 0 {
		fmt.Println("call failed")
	}

	// Unmarshal
	info := &DeviceInfo{}
	if err = json.Unmarshal(resp, info); err != nil {
		fmt.Printf("error: " + err.Error() + conn.Channel.RemoteAddr().String())
		return false
	}
	conn.ID = info.Tid
	info.Conn = conn

	if name, ok := deviceAlias[info.Model]; ok {
		info.Model = name
	}

	if strings.Contains(info.Model, "VMware") {
		info.Model = "VMware"
		info.Model += " " + info.OS
	}
	if strings.Contains(info.Model, "Cloud") {
		info.Model += " " + info.OS
	}

	// add
	ConnectionPool[conn.ID] = info

	fmt.Printf("--> device online: %s\n", info.Model)
	return true
}
