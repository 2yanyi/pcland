package methods

import (
	"r/client/methods/syscontrol"
	"time"
)

const SignRestart = "client.restart"

func Restart(data []byte) ([]byte, error) {
	return nil, syscontrol.Restart()
}

const SignShutdown = "client.shutdown"

func Shutdown(data []byte) ([]byte, error) {
	__work := func() {
		time.Sleep(time.Second)
		syscontrol.Shutdown()
	}
	go __work()
	return nil, nil
}
