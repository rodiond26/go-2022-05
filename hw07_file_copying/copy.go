package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
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
	defer func() {
		closeErr := src.Close()
		err = fmt.Errorf("main error [%s], file close error [%s]", err, closeErr)
	}()

	result, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := result.Close()
		err = fmt.Errorf("main error [%s], file close error [%s]", err, closeErr)
	}()

	_, err = src.Seek(offset, io.SeekStart)

	if err != nil {
		return err
	}

	size := copyLength(fromPath, offset, limit)

	bar := pb.StartNew(int(size))
	for i := 0; i < int(size); i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	defer bar.Finish()

	if limit == 0 {
		written, _ := io.Copy(result, src)
		time.Sleep(time.Millisecond)
		fmt.Printf("written %d\n", written)
	} else {
		written, _ := io.CopyN(result, src, size)
		time.Sleep(time.Millisecond)
		fmt.Printf("written %d\n", written)
	}

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
