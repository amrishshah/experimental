package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct {
	value int32
}

func (c *Counter) Increment() {
	for {
		oldValue := atomic.LoadInt32(&c.value)
		newValue := oldValue + 1
		if atomic.CompareAndSwapInt32(&c.value, oldValue, newValue) {
			return // Exit loop if swap is successful
		}
	}
}

func (c *Counter) GetValue() int32 {
	return atomic.LoadInt32(&c.value)
}

func main() {
	counter := Counter{}

	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Println(counter.GetValue())
}
