package main

import (
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
	if err := tc.Connect(); err != nil {
		log.Println(err)
		return
	}

	defer tc.Close()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	defer signal.Stop(sigCh)

	resultCh := make(chan error)
	run := func(f func() error) {
		resultCh <- f()
	}
	go run(tc.Send)
	go run(tc.Receive)

	select {
	case <-sigCh:
	case err := <-resultCh:
		if err != nil {
			log.Println(err)
		}
	}
}
