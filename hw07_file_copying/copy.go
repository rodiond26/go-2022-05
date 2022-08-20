package main

import (
	"bufio"
	"errors"
	"fmt"
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

	bufferSize := 100
	buffer := make([]byte, bufferSize)
	reader := bufio.NewReader(src)
	writer := bufio.NewWriter(result)

	for {
		_, err := reader.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		written, err := writer.Write(buffer)
		fmt.Printf("written %d\n", written)
		fmt.Printf("buffer %v\n", buffer)
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
	if offset < fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	return nil
}
