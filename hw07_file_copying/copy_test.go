package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	src = "./testdata/input.txt"
	dst = "./testdata/result.txt"
)

type testCase struct {
	src          string
	dst          string
	offset       int64
	limit        int64
	expectedSize int64
}

var tests = []testCase{
	{
		src:          src,
		dst:          dst,
		offset:       0,
		limit:        0,
		expectedSize: 6617,
	},
	{
		src:          src,
		dst:          dst,
		offset:       0,
		limit:        1000,
		expectedSize: 1000,
	},
	{
		src:          src,
		dst:          dst,
		offset:       1000,
		limit:        0,
		expectedSize: 5617,
	},
	{
		src:          src,
		dst:          dst,
		offset:       1000,
		limit:        1000,
		expectedSize: 1000,
	},
}

func TestCopy(t *testing.T) {
	for _, testCase := range tests {
		assert.True(t, !isExists(dst))
		_ = Copy(testCase.src, testCase.dst, testCase.offset, testCase.limit)
		assert.FileExists(t, dst)
		dstInfo, _ := os.Stat(dst)
		assert.Equal(t, testCase.expectedSize, dstInfo.Size())
	}
	defer os.Remove(dst)

	t.Run("offset > size", func(t *testing.T) {
		var offset int64 = 10000
		var limit int64
		require.Equal(t, ErrOffsetExceedsFileSize, Copy(src, dst, offset, limit))
	})

	t.Run("no src", func(t *testing.T) {
		src1 := "testdata/test.txt"
		var offsetFile int64 = 6000
		var limitFile int64
		err := Copy(src1, dst, offsetFile, limitFile)
		require.True(t, errors.Is(err, os.ErrNotExist))
	})

	t.Run("invalid argument", func(t *testing.T) {
		var offsetFile int64 = 6000
		var limitFile int64 = -1
		err := Copy(src, dst, offsetFile, limitFile)
		require.True(t, errors.Is(err, ErrNotValidArgument))
	})
}

func isExists(file string) bool {
	_, err := os.Stat(file)
	return os.IsExist(err)
}
