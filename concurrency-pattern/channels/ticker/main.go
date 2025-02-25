package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)
	var ch = make(chan bool)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Tick at", t)

			case <-ch:
				return

			}
		}
	}()

	time.Sleep(10 * time.Second)
	ch <- true
	fmt.Println("Ticker stopped")
}
