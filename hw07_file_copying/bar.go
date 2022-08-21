package main

import "fmt"

func printProgressBar(written, toWrite int64) {
	symbol := "#"
	percent := int64(float32(written) / float32(toWrite) * 100)
	progressBar := ""
	step := 5

	for i := 0; i < int(percent)/step; i++ {
		progressBar += symbol
	}
	fmt.Printf("\r[%-20s]\t%d%%\t%d/%d\t", progressBar, percent, written, toWrite)
}
