package main

import (
	"errors"
	"os"
	"os/exec"
)

const (
	ExitCodeIOError         = 5
	ExitCodeCommandNotFound = 127
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return ExitCodeCommandNotFound
	}
	command, args := cmd[0], cmd[1:]
	execCmd := exec.Command(command, args...)
	execCmd.Stdin, execCmd.Stdout, execCmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	execCmd.Env = env.Strings()
	if err := execCmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return ExitCodeIOError
	}
	return 0
}
