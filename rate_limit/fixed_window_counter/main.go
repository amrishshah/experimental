package main

import (
	"fmt"
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	mu        sync.Mutex
	limit     int
	counter   int
	resetTime time.Time
	window    time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:     limit,
		window:    window,
		resetTime: time.Now().Add(window),
	}
}

// Allow checks if a request can be processed
func (fw *FixedWindowLimiter) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()

	if now.After(fw.resetTime) {
		fw.counter = 0
		fw.resetTime = now.Add(fw.window)
	}

	if fw.counter < fw.limit {
		fw.counter++
		return true
	}
	return false
}

func main() {
	fw := NewFixedWindowLimiter(1, 2*time.Second)

	for i := 1; i <= 10; i++ {
		if fw.Allow() {
			fmt.Println("Request allowed", i)
		} else {
			fmt.Println("Request denied", i)
		}
		time.Sleep(500 * time.Millisecond) // Simulate request interva
	}
}
