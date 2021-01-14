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

	var err error
	tc := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err = tc.Connect(); err != nil {
		log.Println(err)
		return
	}

	defer tc.Close()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	defer signal.Stop(sigCh)

	doneCh := make(chan struct{})
	run := func(f func() error) {
		defer func() {
			select {
			case doneCh <- struct{}{}:
			default:
			}
		}()
		err = f()
	}
	go run(tc.Send)
	go run(tc.Receive)

	select {
	case <-sigCh:
	case <-doneCh:
	}
	if err != nil {
		log.Println(err)
	}
}
