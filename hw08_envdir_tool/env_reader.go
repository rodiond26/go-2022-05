package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dirPath string) (Environment, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	args := make(Environment)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if err = checkName(entry.Name()); err != nil {
			args[entry.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}
		if err = checkSize(entry); err != nil {
			args[entry.Name()] = EnvValue{Value: "", NeedRemove: true}
		}

		rawArg, err := nextRawArg(dirPath, entry.Name())
		if err != nil {
			return nil, err
		}

		arg := cleanRawArg(rawArg)
		args[entry.Name()] = EnvValue{Value: arg, NeedRemove: false}
	}

	return args, nil
}

func checkName(name string) error {
	if strings.Contains(name, "=") {
		return fmt.Errorf("incorrect file name [%s]", name)
	}
	return nil
}

func checkSize(entry os.DirEntry) error {
	if fileInfo, err := entry.Info(); err != nil {
		return err
	} else if fileInfo.Size() == 0 {
		return fmt.Errorf("incorrect file size [%s]", entry.Name())
	} else {
		return nil
	}
}

func nextRawArg(dirPath, fileName string) (string, error) {
	file, err := os.Open(filepath.Join(dirPath, fileName))
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := file.Close()
		err = fmt.Errorf("main error [%s], file [%s] close error [%s]", err, fileName, closeErr)
	}()

	rawArg, _, err := bufio.NewReader(file).ReadLine()
	if err != nil {
		return "", err
	}
	return string(rawArg), nil
}

func cleanRawArg(arg string) string {
	if strings.IndexByte(arg, '\x00') != -1 {
		arg = strings.ReplaceAll(arg, "\x00", "\n")
	}
	arg = strings.TrimRight(arg, " ")
	arg = strings.TrimRight(arg, "\n")
	return arg
}
