package main

import (
	"os"
	"os/exec"
)

const (
	successCode            = 0
	unsuccessfulChildCode  = 1
	unsuccessfulEnvdirCode = 111
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	childCmd := cmd[0]
	childArgs := cmd[1:]

	command := exec.Command(childCmd, childArgs...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for key, envValue := range env {
		// setting env vars
		if envValue.NeedRemove {
			if _, ok := os.LookupEnv(key); ok {
				_ = os.Unsetenv(key)
			}
			continue
		}
		if _, ok := os.LookupEnv(key); ok {
			_ = os.Unsetenv(key)
		}

		err := os.Setenv(key, envValue.Value)
		if err != nil {
			return unsuccessfulEnvdirCode
		}
	}

	command.Env = os.Environ()
	if err := command.Run(); err != nil {
		return unsuccessfulChildCode
	}

	return successCode
}
