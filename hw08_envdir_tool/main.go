package main

import "os"

func main() {
	args := os.Args
	dirPath := args[1]
	childCmdWithArgs := args[2:]

	envValues, err := ReadDir(dirPath)
	if err != nil {
		os.Exit(unsuccessfulEnvdirCode)
	}
	exitCode := RunCmd(childCmdWithArgs, envValues)
	os.Exit(exitCode)
}
