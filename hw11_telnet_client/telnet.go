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
	tcp = "tcp"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &tnClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type tnClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (client *tnClient) Connect() error {
	conn, err := net.DialTimeout(tcp, client.address, client.timeout)
	if err != nil {
		return fmt.Errorf("When connect to [%v] then error [%v]", client.address, err)
	}
	client.conn = conn
	return nil
}

func (t *tnClient) Close() error {
	return t.conn.Close()
}

func (client *tnClient) Send() error {
	scanner := bufio.NewScanner(client.in)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := client.conn.Write([]byte(fmt.Sprintln(text)))
		if err != nil {
			return fmt.Errorf("...Connection was closed by peer")
		}
	}

	fmt.Fprintln(os.Stderr, "...EOF")
	return scanner.Err()
}

func (client *tnClient) Receive() error {
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintln(client.out, text)
	}

	return scanner.Err()
}
