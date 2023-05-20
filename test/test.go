package main

import (
	"bufio"
	"fmt"
	"os"
	"r/pkg/cmdline"
)

var reader = bufio.NewReader(os.Stdin)

func main() {
	go cmdline.Accept("C:\\")
	go func() {
		for line := range cmdline.OUTPUT {
			fmt.Printf("%s", line)
		}
	}()

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		select {
		case cmdline.INPUT <- string(line):
		default:
			return
		}
	}
}
