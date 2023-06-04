package methods

import (
	"github.com/google/uuid"
	"os"
	"r/client/sysinfo"
	"r/device"
	"r2/pkg/generic2/chars"
	"runtime"
	"time"
)

func GetTerminalID() string {
	idFp := "/tmp/PCland-TerminalID.txt"
	if runtime.GOOS == "windows" {
		idFp = os.Getenv("TMP") + "\\PCland-TerminalID.txt"
	}

again:
	tid := chars.CatString(idFp)
	if tid == "" {
		if err := os.WriteFile(idFp, []byte(uuid.New().String()), 0666); err != nil {
			panic(err)
		}
		goto again
	}

	return tid
}

const SignGetDeviceInfo = "client.deviceInfo"

func GetDeviceInfo(data []byte) ([]byte, error) {

	target, title := sysinfo.GetTarget()

	info := &device.DeviceInfo{
		Power:  int(CPUPower()),
		Target: target,
		OS:     title,
		Model:  sysinfo.GetProductName(),
		Addr:   sysinfo.GetLanIPv4(),
		Tid:    GetTerminalID(),
	}

	return chars.ToJsonBytes(info, ""), nil
}

func CPUPower() (count int32) {
	go func() {
		defer func() {
			println("POWER", count)
			if err := recover(); err != nil {
			}
		}()
		off = false
		for i := 2; ; i++ {
			if fibonacci(i, &count) == 0 {
				return
			}
		}
	}()
	time.Sleep(time.Second)
	off = true

	return
}

func fibonacci(n int, count *int32) int {
	*count++

	if off {
		return 0
	}
	if n < 2 {
		return 1
	}

	return fibonacci(n-1, count) + fibonacci(n-2, count)
}

var off bool
