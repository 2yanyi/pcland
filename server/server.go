// Source code file, created by Developer@YANYINGSONG.

package server

import (
	"encoding/json"
	"fmt"
	"library/console"
	"library/generic/chars"
	"library/generic/errcause"
	"r/client/methods"
	"r/device"
	"rpcpd"
	"rpcpd/seqx"
	"strings"
	"time"
)

const ListenPort = ":11299"

var Power *rpcpd.Power

var deviceAlias = make(map[string]string)

func InitDeviceAlias() {
	data := chars.CatBytes("alias.json")
	if err := json.Unmarshal(data, &deviceAlias); err != nil {
		println(err.Error())
	}
}

func deleteConnection(conn *rpcpd.Conn) {
	console.INFO("disconnect %s", conn.Channel.RemoteAddr())
	conn.Channel.Close()
	delete(device.Pool, conn.ID)
}

func addConnection(conn *rpcpd.Conn) bool {
	console.INFO("connect %s", conn.Channel.RemoteAddr().String())
	time.Sleep(3 * time.Second)

	// call client get device info.
	resp, err := Power.Call(conn, []byte(methods.SignGetDeviceInfo), nil)
	if err != nil {
		console.ERROR(fmt.Errorf("failed call: %s", err))
		return false
	}
	if resp == nil {
		console.ERROR(fmt.Errorf("failed call: resp is nil"))
	}

	// Unmarshal
	info := &device.DeviceInfo{}
	if err = json.Unmarshal(resp, info); err != nil {
		console.ERROR(fmt.Errorf("unmarshal DeviceInfo{} %s", resp))
		return false
	}
	conn.ID = info.Tid
	info.Conn = conn

	if name, ok := deviceAlias[info.Model]; ok {
		info.Model = name
	}

	if strings.Contains(info.Model, "VMware") {
		info.Model = "vmw_" + info.OS
	}
	if strings.Contains(info.Model, "Cloud") {
		info.Model += " " + info.OS
	}

	// Add
	device.Pool[conn.ID] = info
	console.INFO("Device online: %s", chars.ToJsonBytes(info, ""))

	return true
}

func RunServer() error {
	seqx.Init(2)

	var powerErr error
	Power, powerErr = rpcpd.Listen(ListenPort)
	if powerErr != nil {
		panic(powerErr)
	}
	defer Power.Listener.Close()

	for {
		conn, err := Power.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		go connectionAdd(conn)
		go connectionLive(conn)
		go connectionPing(conn)
	}
}

func connectionAdd(conn *rpcpd.Conn) {
	defer errcause.Recover()
	if !addConnection(conn) {
		conn.Channel.Close()
	}
}

func connectionLive(conn *rpcpd.Conn) {
	defer errcause.Recover()
	if err := Power.ConnectionProcessor(conn); err != nil {
		console.WARN(err.Error())
	}
}

func connectionPing(conn *rpcpd.Conn) {
	defer errcause.Recover()
	defer deleteConnection(conn)

	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for range t.C {
		pong, _ := Power.Call(conn, []byte(methods.SignPing), nil)
		if string(pong) != methods.SignPingRET {
			return
		}
	}
}
