// Source code file, created by Developer@YANYINGSONG.

package cpower

import "time"

var off bool

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
