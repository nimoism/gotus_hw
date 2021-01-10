package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	logPrefix    = "..."
	logEOF       = "EOF"
	logClosed    = "Connection was closed by peer"
	logConnected = "Connected to %s"
)

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.Reader
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.Reader, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() (err error) {
	if c.conn, err = net.DialTimeout("tcp", c.address, c.timeout); err != nil {
		return err
	}
	return c.log(logConnected, c.address)
}

func (c *telnetClient) Send() error {
	return c.handleMessages(c.in, c.conn, logEOF)
}

func (c *telnetClient) Receive() error {
	return c.handleMessages(c.conn, c.out, logClosed)
}

func (c *telnetClient) handleMessages(in io.Reader, out io.Writer, eofMsg string) error {
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return c.log(eofMsg)
}

func (c *telnetClient) Close() error {
	return c.conn.Close()
}

func (c *telnetClient) log(s string, v ...interface{}) error {
	_, err := fmt.Fprintf(os.Stderr, "%s%s\n", logPrefix, fmt.Sprintf(s, v...))
	return err
}
