package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 1; j <= 3; j++ {
				ch <- id*10 + j // Send data
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// Receive data from the channel
	for val := range ch {
		fmt.Println(val)
	}

}
