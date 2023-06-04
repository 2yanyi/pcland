// Source code file, created by Developer@YANYINGSONG.

package syscontrol

func Shutdown() error {
	_ = hibernate()
	return shutdown()
}

func Restart() error {
	return restart()
}
