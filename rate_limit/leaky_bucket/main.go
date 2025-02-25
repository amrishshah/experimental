package main

import (
	"fmt"
	"sync"
	"time"
)

// LeakyBucket struct
type LeakyBucket struct {
	mu       sync.Mutex
	capacity int
	rate     time.Duration
	queue    []time.Time
}

// NewLeakyBucket initializes a new leaky bucket
func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		rate:     rate,
		queue:    []time.Time{},
	}
}

// Allow checks if a request can be processed
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()

	// Remove expired requests
	for len(lb.queue) > 0 && now.Sub(lb.queue[0]) > lb.rate {
		lb.queue = lb.queue[1:]
	}

	// Allow request if bucket is not full
	if len(lb.queue) < lb.capacity {
		lb.queue = append(lb.queue, now)
		return true
	}

	return false
}

func main() {
	limiter := NewLeakyBucket(5, 500*time.Millisecond) // 5 requests allowed, leaking every 500ms

	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Println("Request allowed", i)
		} else {

			fmt.Println("Request denied", i)
		}
		time.Sleep(200 * time.Millisecond) // Simulate request interval
	}
}
