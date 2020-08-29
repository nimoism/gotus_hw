package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	fmt.Println("current time:", time.Now().Round(time.Second))
	exactTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		log.Fatalf("Getting time from ntp server failed: %v", err)
	}
	fmt.Println("exact time:", exactTime.Round(time.Second))
}
