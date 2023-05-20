// Source code file, created by Developer@YANYINGSONG.

package syscontrol

import (
	"os/exec"
)

func hibernate() error {
	return exec.Command("rundll32.exe", "powrprof.dll", "SetSuspendState").Run()
}

func shutdown() error {
	return exec.Command("shutdown", "/s", "/t", "0").Run()
}

func restart() error {
	return exec.Command("shutdown", "/r", "/t", "0").Run()
}
