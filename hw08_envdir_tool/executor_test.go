package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	childCmd  = "testdata/echo.sh"
	childArgs = []string{"arg1", "arg2"}
)

func TestRunCmd(t *testing.T) {
	t.Run("success code", func(t *testing.T) {
		command := []string{childCmd}
		command = append(command, childArgs...)

		args := make(Environment)
		args["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		args["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		args["FOO"] = EnvValue{Value: "foo", NeedRemove: false}
		args["ADDED"] = EnvValue{Value: "added", NeedRemove: false}
		args["UNSET"] = EnvValue{Value: "", NeedRemove: true}

		actualCode := RunCmd(command, args)
		expectedCode := successCode

		require.Equal(t, expectedCode, actualCode)
	})

	t.Run("no args", func(t *testing.T) {
		command := []string{childCmd}
		command = append(command, childArgs...)

		args := make(Environment)
		actualCode := RunCmd(command, args)
		expectedCode := successCode

		require.Equal(t, expectedCode, actualCode)
	})

	t.Run("nil args", func(t *testing.T) {
		command := []string{childCmd}
		command = append(command, childArgs...)

		actualCode := RunCmd(command, nil)
		expectedCode := successCode

		require.Equal(t, expectedCode, actualCode)
	})

	t.Run("empty command", func(t *testing.T) {
		command := make([]string, 10)

		args := make(Environment)
		args["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		args["EMPTY"] = EnvValue{Value: "", NeedRemove: false}

		actualCode := RunCmd(command, args)
		expectedCode := unsuccessfulChildCode

		require.Equal(t, expectedCode, actualCode)
	})

	t.Run("empty command, nil args", func(t *testing.T) {
		command := make([]string, 10)

		actualCode := RunCmd(command, nil)
		expectedCode := unsuccessfulChildCode

		require.Equal(t, expectedCode, actualCode)
	})

	t.Run("incorrect child command path", func(t *testing.T) {
		wrongChildCmd := "/echo.sh"
		command := []string{wrongChildCmd}
		command = append(command, childArgs...)

		args := make(Environment)
		args["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		args["EMPTY"] = EnvValue{Value: "", NeedRemove: true}

		actualCode := RunCmd(command, args)
		expectedCode := unsuccessfulChildCode

		require.Equal(t, expectedCode, actualCode)
	})

	t.Run("incorrect arg", func(t *testing.T) {
		command := []string{childCmd}
		command = append(command, childArgs...)

		args := make(Environment)
		args["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		args["EMPTY"] = EnvValue{Value: "", NeedRemove: true}
		args["SOME_VAR"] = EnvValue{Value: "a\x00 ", NeedRemove: false}

		actualCode := RunCmd(command, args)
		expectedCode := unsuccessfulEnvdirCode

		require.Equal(t, expectedCode, actualCode)
	})
}
