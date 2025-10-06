package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type ConcurrentQueue struct {
	queue []int32
	mu    sync.Mutex
}

func (q *ConcurrentQueue) Enqueue(item int32) {
	q.mu.Lock()
	q.queue = append(q.queue, item)
	q.mu.Unlock()
}

func (q *ConcurrentQueue) Dequeue() int32 {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		panic("No item present")
	}
	item := q.queue[0]
	q.queue = q.queue[1:]
	return item
}

var wgE sync.WaitGroup
var wgD sync.WaitGroup

func main() {
	q := ConcurrentQueue{
		queue: make([]int32, 0),
		mu:    sync.Mutex{},
	}

	for i := 0; i < 1000000; i++ {
		wgE.Add(1)
		go func() {
			q.Enqueue(rand.Int31())
			defer wgE.Done()
		}()
	}

	for i := 0; i < 1000000; i++ {
		wgD.Add(1)
		go func() {
			q.Dequeue()
			defer wgD.Done()
		}()
	}

	wgE.Wait()
	wgD.Wait()
	fmt.Println(len(q.queue))

	// fmt.Println(q.Dequeue())
	// fmt.Println(q.Dequeue())
	// fmt.Println(q.Dequeue())
	// fmt.Println(q.Dequeue())

}
