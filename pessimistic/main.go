package main

import (
	"fmt"
	"sync"
)

var count int

func inc(wg *sync.WaitGroup, s *sync.Mutex) {
	s.Lock()
	count++
	s.Unlock()
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	var s sync.Mutex

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go inc(&wg, &s)
	}
	wg.Wait()
	fmt.Println(count)
}
