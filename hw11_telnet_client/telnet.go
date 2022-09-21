package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	tcp        = "tcp"
	connClosed = "...Connection was closed by peer"
	eof        = "...EOF"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (client *telnetClient) Connect() error {
	conn, err := net.DialTimeout(tcp, client.address, client.timeout)
	if err != nil {
		return fmt.Errorf("when connect to [%s] then error [%w]", client.address, err)
	}
	client.conn = conn
	return nil
}

func (client *telnetClient) Close() error {
	return client.conn.Close()
}

func (client *telnetClient) Send() error {
	scanner := bufio.NewScanner(client.in)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := client.conn.Write([]byte(fmt.Sprintln(text)))
		if err != nil {
			return fmt.Errorf(connClosed)
		}
	}

	fmt.Fprintln(os.Stderr, eof)
	return scanner.Err()
}

func (client *telnetClient) Receive() error {
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintln(client.out, text)
	}

	return scanner.Err()
}
