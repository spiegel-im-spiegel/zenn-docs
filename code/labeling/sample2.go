//go:build run
// +build run

package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	abort := make(chan os.Signal, 1)
	defer func() {
		close(abort)
		fmt.Println("Close abort channel.")
	}()
	signal.Notify(abort, os.Interrupt)

	fmt.Println("Commencing countdown.  Press Ctrl+C to abort.")
	tick := time.Tick(1 * time.Second)
LOOP:
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick:
			//何もしない
		case <-abort:
			fmt.Println("Countdown aborted!")
			break LOOP
		}
	}
}
