package cmdline

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"r2/pkg/generic2/chars/cat"
	"strings"
)

var INPUT = make(chan string)
var OUTPUT = make(chan string, 1024)

func display(s string) {
	OUTPUT <- s + "\n"
}

func displayError(err error) {
	display("error: " + err.Error())
}

func execute(input string) {

	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", input)
	cmd.Dir = workDir

	errPipe, err1 := cmd.StderrPipe()
	outPipe, err2 := cmd.StdoutPipe()
	if err1 != nil && err2 != nil {
		displayError(err1)
		displayError(err2)
		return
	}
	if err := cmd.Start(); err != nil {
		displayError(err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(errPipe)
		for scanner.Scan() {
			display("err | " + scanner.Text())
		}
	}()
	scanner := bufio.NewScanner(outPipe)
	for scanner.Scan() {
		display(scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		displayError(errors.New("command execution failed"))
	}

	return
}

func connection(input string) {
	defer func() {
		OUTPUT <- fmt.Sprintf("BASH %s# ", workDir)
	}()

	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	if strings.HasPrefix(input, "cd ") {
		dir := strings.TrimSpace(input[3:])
		input = ""
		if !cat.FileExist(dir) {
			displayError(fmt.Errorf("chdir %s: The system cannot find the path specified", dir))
			return
		}
		workDir = dir
	}

	// 执行命令
	execute(input)
}

func Accept(dir string) {
	workDir = dir
	OUTPUT <- fmt.Sprintf("BASH %s# ", workDir)

	for input := range INPUT {
		if input == "exit" {
			return
		}

		connection(input)
	}
}

var workDir string
