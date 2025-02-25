package main

import (
	"sync"
	"sync/atomic"
)

var (
	count int32 = 0
	mu    sync.Mutex
)

func updateCounter(count *int32, wg *sync.WaitGroup) {
	// mu.Lock()
	// count++
	// mu.Unlock()
	atomic.AddInt32(count, 1)
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 500000; i++ {
		wg.Add(1)
		go updateCounter(&count, &wg)
	}

	wg.Wait()

	println(count)

}
