package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("check vars", func(t *testing.T) {
		rootPath, _ := os.Getwd()
		dirPath := "/testdata/env"
		path := filepath.Join(rootPath, dirPath)

		actualEnv, actualErr := ReadDir(path)
		expectedEnv := make(Environment)
		expectedEnv["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		expectedEnv["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		expectedEnv["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		expectedEnv["HELLO"] = EnvValue{Value: `"hello"`, NeedRemove: false}
		expectedEnv["UNSET"] = EnvValue{Value: "", NeedRemove: true}

		require.Equal(t, expectedEnv, actualEnv)
		require.Nil(t, actualErr)
	})

	t.Run("wrong directory", func(t *testing.T) {
		rootPath, _ := os.Getwd()
		wrongDirPath := "/test"
		path := filepath.Join(rootPath, wrongDirPath)

		actualEnv, actualErr := ReadDir(path)
		require.Nil(t, actualEnv)
		require.True(t, errors.Is(actualErr, os.ErrNotExist))
	})
}
