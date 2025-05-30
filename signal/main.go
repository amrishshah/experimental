package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(2)
	go temp(&wg)
	go waitSignal(&wg, ch)
	wg.Wait()
}

func temp(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("we are in temp")
}

func waitSignal(wg *sync.WaitGroup, ch chan os.Signal) {
	defer wg.Done()
	<-ch
	fmt.Println("Received signal, exiting...")
}
