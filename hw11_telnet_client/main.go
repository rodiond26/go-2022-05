package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

const defaultTimeout = 10

var (
	timeout         time.Duration
	invalidArgs     = "invalid number of args"
	ErrConnection   = errors.New("connection error")
	closeConnection = "close connection error"
)

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*defaultTimeout, "connection timeout [default=10s]")
}

func main() {
	flag.Parse()
	flagArgs := flag.Args()
	if len(flagArgs) != 2 {
		fmt.Println(invalidArgs)
		os.Exit(1)
	}

	host, port := flagArgs[0], flagArgs[1]
	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := run(client); err != nil {
		os.Exit(1)
	}
}

func run(client TelnetClient) error {
	if err := client.Connect(); err != nil {
		return ErrConnection
	}
	defer func(client TelnetClient) {
		err := client.Close()
		if err != nil {
			fmt.Println(closeConnection)
		}
	}(client)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		defer stop()
		if err := client.Send(); err != nil {
			log.Fatalf("Error reading from channel. Error: %s", err)
		}
	}()

	go func() {
		defer stop()
		if err := client.Receive(); err != nil {
			log.Fatalf("Error reading from channel. Error: %s", err)
		}
	}()

	<-ctx.Done()
	return nil
}
