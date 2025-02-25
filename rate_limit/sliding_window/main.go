package main

import (
	"fmt"
	"sync"
	"time"
)

// SlidingWindowLimiter struct
type SlidingWindowLimiter struct {
	mu     sync.Mutex
	window time.Duration
	rate   int
	times  []time.Time
}

// NewSlidingWindowLimiter initializes a new sliding window rate limiter
func NewSlidingWindowLimiter(rate int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		rate:   rate,
		window: window,
		times:  []time.Time{},
	}
}

// Allow checks if a request can be processed
func (sw *SlidingWindowLimiter) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	// Remove expired timestamps
	for len(sw.times) > 0 && now.Sub(sw.times[0]) > sw.window {
		sw.times = sw.times[1:]
	}

	// Allow request if within rate limit
	if len(sw.times) < sw.rate {
		sw.times = append(sw.times, now)
		return true
	}

	return false
}

// Example Usage
func main() {
	limiter := NewSlidingWindowLimiter(2, 2*time.Second) // 5 requests in 2 seconds

	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Println("Request allowed", i)
		} else {
			fmt.Println("Request denied", i)
		}
		time.Sleep(400 * time.Millisecond) // Simulate request interval
	}
}
