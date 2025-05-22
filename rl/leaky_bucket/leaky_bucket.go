package leaky_bucket

import (
	"math"
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity     int           // Maximum bucket capacity
	leakRate     time.Duration // Rate at which tokens leak (time per request)
	tokens       int           // Current number of tokens in bucket
	lastLeakTime time.Time     // Last time tokens were leaked
	mutex        sync.Mutex    // For thread safety
}

func NewLeakyBucket(capacity int, requestsPerSecond int) *LeakyBucket {
	var leakRate time.Duration
	if requestsPerSecond <= 0 {
		// Handle zero or negative rate by setting a very slow leak rate
		// or using a special value to indicate "no leaking"
		leakRate = time.Duration(math.MaxInt64) // Effectively no leaking
	} else {
		leakRate = time.Second / time.Duration(requestsPerSecond)
	}

	return &LeakyBucket{
		capacity:     capacity,
		leakRate:     leakRate,
		tokens:       0,
		lastLeakTime: time.Now(),
	}
}

// leak removes tokens based on elapsed time
func (lb *LeakyBucket) leak() {
	now := time.Now()
	elapsed := now.Sub(lb.lastLeakTime)

	// Calculate how many tokens should have leaked since last check
	leakedTokens := int(elapsed / lb.leakRate)

	if leakedTokens > 0 {
		// Update tokens and last leak time
		if leakedTokens > lb.tokens {
			lb.tokens = 0
		} else {
			lb.tokens -= leakedTokens
		}

		// Update the last leak time based on actual tokens leaked
		lb.lastLeakTime = lb.lastLeakTime.Add(time.Duration(leakedTokens) * lb.leakRate)
	}
}

// Allow checks if a request can be processed
func (lb *LeakyBucket) Allow() bool {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	// Leak tokens based on elapsed time
	lb.leak()

	// If the bucket is full, reject the request
	if lb.tokens >= lb.capacity {
		return false
	}

	// Add token to bucket and allow request
	lb.tokens++
	return true
}

// Take blocks until a request can be processed
func (lb *LeakyBucket) Take() time.Time {
	for {
		lb.mutex.Lock()
		lb.leak()

		if lb.tokens < lb.capacity {
			lb.tokens++
			now := time.Now()
			lb.mutex.Unlock()
			return now
		}

		// Calculate time to wait for the next token to leak
		waitTime := lb.leakRate
		lb.mutex.Unlock()

		time.Sleep(waitTime)
	}
}
