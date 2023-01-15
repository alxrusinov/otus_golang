package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, currentEnv := range env {
		if currentEnv.NeedRemove {
			_ = os.Unsetenv(key)
			continue
		}

		if _, ok := os.LookupEnv(key); ok {
			_ = os.Unsetenv(key)
		}

		_ = os.Setenv(key, currentEnv.Value)
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
