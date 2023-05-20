// Source code file, created by Developer@YANYINGSONG.

package connect

import (
	"github.com/google/uuid"
	"os"
	"r/pkg/cpower"
	"r/pkg/syscontrol"
	"r/pkg/sysinfo"
	"r2/pkg/generic2/chars"
	"r2/pkg/generic2/chars/cat"
	"runtime"
	"time"
)

func GetTerminalID() string {
	idFp := "/tmp/PCland-TerminalID.txt"
	if runtime.GOOS == "windows" {
		idFp = os.Getenv("TMP") + "\\PCland-TerminalID.txt"
	}

again:
	tid := cat.String(idFp)
	if tid == "" {
		if err := os.WriteFile(idFp, []byte(uuid.New().String()), 0666); err != nil {
			panic(err)
		}
		goto again
	}

	return tid
}

const FunctionSignGetDeviceInfo = "client.deviceInfo"

func GetDeviceInfo(data []byte) ([]byte, error) {

	target, title := sysinfo.GetTarget()

	info := &DeviceInfo{
		Power:  int(cpower.CPUPower()),
		Target: target,
		OS:     title,
		Model:  sysinfo.GetProductName(),
		Addr:   sysinfo.GetLanIPv4(),
		Tid:    GetTerminalID(),
	}

	return chars.ToJsonBytes(info, ""), nil
}

const FunctionSignRestart = "client.restart"

func Restart(data []byte) ([]byte, error) {
	return nil, syscontrol.Restart()
}

const FunctionSignShutdown = "client.shutdown"

func Shutdown(data []byte) ([]byte, error) {
	__work := func() {
		time.Sleep(time.Second)
		syscontrol.Shutdown()
	}
	go __work()
	return nil, nil
}
