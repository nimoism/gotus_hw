package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "")
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Params `host` and `port` are mandatory")
	}
	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	tc := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	logErr(tc.Connect())
	defer func() { logErr(tc.Close()) }()

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGINT)
	defer signal.Stop(sigCh)

	run := func(f func() error) {
		defer cancel()
		logErr(f())
	}
	go run(tc.Send)
	go run(tc.Receive)

	select {
	case <-sigCh:
		cancel()
	case <-ctx.Done():
	}
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
