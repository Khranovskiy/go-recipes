package main

import (
	"fmt"
	"time"

	"go.uber.org/ratelimit"

	"github.com/Khranovskiy/go-recipes/rl/fixed_window"
	"github.com/Khranovskiy/go-recipes/rl/leaky_bucket"
	"github.com/Khranovskiy/go-recipes/rl/token_bucket"
)

// leaky bucket rate limiter

func mainLB() {
	const (
		capacity          = 5
		requestsPerSecond = 2
	)
	limiter := leaky_bucket.NewLeakyBucket(capacity, requestsPerSecond)

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := limiter.Take()
		fmt.Printf("Request %d: waited %v\n", i, now.Sub(prev))
		prev = now
	}
}

func mainTB() {
	const (
		maxTokens        = 5
		refillRatePerSec = 1
	)
	rateLimiter := token_bucket.NewTokenBucket(maxTokens, refillRatePerSec)

	// Simulate 20 requests with 500ms between each
	for i := 0; i < 20; i++ {
		allowed := rateLimiter.Allow(1)
		fmt.Printf("Request %d: %v\n", i+1, allowed)
		time.Sleep(500 * time.Millisecond)
	}
}

func mainThrottled() {
	work := func() {
		fmt.Print(".")
	}

	const rps = 1
	handle, cancel := fixed_window.Throttle(rps, work)
	defer cancel()

	start := time.Now()
	const n = 10
	for range n {
		handle()
	}
	fmt.Println()
	fmt.Printf("%d queries took %v\n", n, time.Since(start))
}

func main() {
	mainLB()
	// mainTB()
	// mainThrottled()

	// Create a rate limiter for 100 requests per second
	rl := ratelimit.New(100)

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
}
