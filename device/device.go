package device

import "rpcpd"

var Pool = make(map[string]*DeviceInfo, 10)

type DeviceInfo struct {
	Power  int    `json:"power,omitempty"`
	Target string `json:"target,omitempty"`
	OS     string `json:"os,omitempty"`
	Model  string `json:"model,omitempty"`
	Addr   string `json:"addr,omitempty"`
	Tid    string `json:"tid,omitempty"`

	Conn *rpcpd.Conn `json:"-"`
}
