package program

import (
	"bytes"
	"os/exec"
	"runtime"
)

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
		return stderr, nil
	}

	return stdout, nil
}
