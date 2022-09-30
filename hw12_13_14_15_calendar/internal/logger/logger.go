package logger

import "fmt"

type Logger struct { // TODO
}

func New(level string) *Logger {
	return &Logger{}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Debug(msg string) {
	fmt.Println(msg) // TODO
}

func (l Logger) Warn(msg string) {
	fmt.Println(msg) // TODO
}

func (l Logger) Error(msg string) {
	fmt.Println(msg) // TODO
}
