package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	phrase := "Hello, OTUS!"
	reversedPhrase := stringutil.Reverse(phrase)

	fmt.Println(reversedPhrase)
}
