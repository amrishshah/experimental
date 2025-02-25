package main

import (
	"sync"
)

func temp(msg string, wg *sync.WaitGroup) {
	defer wg.Done()
	println(msg)
	println("hi this is first go routine")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go temp("wewe", &wg)
	wg.Wait()
	println("hi this is end main routine")
}
