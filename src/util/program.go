package program

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

// Program being used
type Program struct {
	Name                 string
	RespectsEndOfOptions bool
}

// Program constructor
func NewProgram(name string, respectsEndOfOptions bool) Program {
	program := Program{
		Name:                 name,
		RespectsEndOfOptions: respectsEndOfOptions,
	}
	return program
}

// Execute the given command in either bash or powershell depending on the detected os
func shellout(command string, silent bool) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := &exec.Cmd{}
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-command", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func Run(command string) (res string, err error) {

	stdout, stderr, err := shellout(command, true)
	if err != nil {
		stderr = stderr + fmt.Sprint(err)
		return stderr, nil
	}

	return stdout, nil
}
