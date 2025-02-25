package main

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter struct
type RateLimiter struct {
	mu        sync.Mutex
	tokens    int
	maxTokens int
	interval  time.Duration
}

// NewRateLimiter initializes a new RateLimiter
func NewRateLimiter(rate int, maxTokens int) *RateLimiter {
	rl := &RateLimiter{
		tokens:    maxTokens,
		maxTokens: maxTokens,
		interval:  time.Second / time.Duration(rate),
	}

	// Background token refill
	go func() {

		ticker := time.NewTicker(rl.interval)
		defer ticker.Stop()
		for range ticker.C {
			println("Tokken refill")
			println(rl.tokens)
			rl.mu.Lock()
			if rl.tokens < rl.maxTokens {
				rl.tokens++
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

// Allow checks if a request can be processed
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	print(rl.tokens)
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// Example Usage
func main() {
	limiter := NewRateLimiter(10, 5) // 2 requests/sec, max burst 5

	for i := 1; i <= 50; i++ {
		if limiter.Allow() {
			fmt.Println("Request allowed", i)
		} else {
			fmt.Println("Request denied", i)
		}
		time.Sleep(100 * time.Millisecond) // Simulate request interval
	}
}
