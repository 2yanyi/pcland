// Source code file, created by Developer@YANYINGSONG.

package syscontrol

import (
	"os/exec"
)

func hibernate() error {
	return exec.Command("systemctl", "suspend").Run()
}

func shutdown() error {
	return exec.Command("shutdown", "-h", "now").Run()
}

func restart() error {
	return exec.Command("reboot").Run()
}
