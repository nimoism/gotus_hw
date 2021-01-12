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

	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	tc := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := tc.Connect(); err != nil {
		return
	}
	defer func() { err = tc.Close() }()

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGINT)
	defer signal.Stop(sigCh)

	run := func(f func() error) {
		defer cancel()
		err = f()
	}
	go run(tc.Send)
	go run(tc.Receive)

	select {
	case <-sigCh:
	case <-ctx.Done():
	}
}
