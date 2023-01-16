package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, currentEnv := range env {
		if _, ok := os.LookupEnv(key); ok {
			_ = os.Unsetenv(key)
		}

		if !currentEnv.NeedRemove {
			_ = os.Setenv(key, currentEnv.Value)
			continue
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	var targetErr *exec.ExitError

	if err := command.Run(); errors.As(err, &targetErr) {
		returnCode = targetErr.ExitCode()
	}

	return
}
