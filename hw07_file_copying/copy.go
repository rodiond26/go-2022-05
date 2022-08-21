package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNotValidArgument      = errors.New("argument(s) is not valid")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := checkArgs(fromPath, toPath, offset, limit)
	if err != nil {
		return err
	}

	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	result, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer result.Close()

	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	toWrite := copyLength(fromPath, offset, limit)
	var buffer int64 = 7
	var written int64 = 0
	var sum int64 = 0

	for toWrite-sum > buffer {
		written, err = io.CopyN(result, src, buffer)
		sum += written
		printProgressBar(sum, toWrite)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
	}

	if toWrite-sum != 0 {
		buffer = toWrite - sum
		written, err := io.CopyN(result, src, buffer)
		sum += written
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
	printProgressBar(sum, toWrite)

	return nil
}

func checkArgs(srcPath, destPath string, offset, limit int64) error {
	if srcPath == "" || destPath == "" || offset < 0 || limit < 0 {
		return ErrNotValidArgument
	}
	if srcPath == destPath && offset == 0 && limit == 0 {
		return errors.New("it's no need to copy the same file")
	}

	fileInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

func copyLength(srcPath string, offset, limit int64) int64 {
	fileInfo, _ := os.Stat(srcPath)

	if limit == 0 {
		return fileInfo.Size() - offset
	} else {
		if offset+limit > fileInfo.Size() {
			return fileInfo.Size() - offset
		} else {
			return limit
		}
	}
}
