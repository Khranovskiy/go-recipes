package leaky_bucket

import (
	"sync"
	"testing"
	"time"
)

func TestLeakyBucketCreation(t *testing.T) {
	capacity := 5
	rps := 10

	lb := NewLeakyBucket(capacity, rps)

	if lb.capacity != capacity {
		t.Errorf("Expected capacity to be %d, got %d", capacity, lb.capacity)
	}

	expectedLeakRate := time.Second / time.Duration(rps)
	if lb.leakRate != expectedLeakRate {
		t.Errorf("Expected leak rate to be %v, got %v", expectedLeakRate, lb.leakRate)
	}

	if lb.tokens != 0 {
		t.Errorf("Expected initial tokens to be 0, got %d", lb.tokens)
	}
}

func TestAllowMethod(t *testing.T) {
	// Create a bucket with capacity 3 and 1 request per second
	lb := NewLeakyBucket(3, 1)

	// First 3 requests should be allowed (filling the bucket)
	for i := 0; i < 3; i++ {
		if !lb.Allow() {
			t.Errorf("Request %d should be allowed", i)
		}
	}

	// 4th request should be rejected (bucket is full)
	if lb.Allow() {
		t.Errorf("4th request should be rejected")
	}
}

func TestLeakingTokens(t *testing.T) {
	// Create a bucket with capacity 5 and 10 requests per second (100ms per request)
	lb := NewLeakyBucket(5, 10)

	// Fill the bucket
	for i := 0; i < 5; i++ {
		lb.Allow()
	}

	// Bucket should be full now
	if lb.Allow() {
		t.Errorf("Bucket should be full, request should be rejected")
	}

	// Wait for 200ms, which should leak 2 tokens
	lb.lastLeakTime = lb.lastLeakTime.Add(-200 * time.Millisecond)

	// Now we should be able to add 2 more tokens
	if !lb.Allow() {
		t.Errorf("Should allow request after leaking")
	}

	if !lb.Allow() {
		t.Errorf("Should allow second request after leaking")
	}

	// But third should fail
	if lb.Allow() {
		t.Errorf("Third request should fail, bucket should be full again")
	}
}

func TestTakeMethod(t *testing.T) {
	// Create a bucket with capacity 3 and 10 requests per second
	lb := NewLeakyBucket(3, 10)

	// First 3 calls to Take() should return immediately
	start := time.Now()
	for i := 0; i < 3; i++ {
		lb.Take()
	}
	elapsed := time.Since(start)

	if elapsed > 50*time.Millisecond {
		t.Errorf("First 3 Take() calls should return almost immediately, took %v", elapsed)
	}

	// Next Take() should block for about 100ms
	start = time.Now()
	done := make(chan bool)

	go func() {
		lb.Take()
		done <- true
	}()

	select {
	case <-done:
		elapsed = time.Since(start)
		// Should wait approximately 100ms (leakRate)
		if elapsed < 90*time.Millisecond || elapsed > 150*time.Millisecond {
			t.Errorf("Take() should block for ~100ms, actually blocked for %v", elapsed)
		}
	case <-time.After(200 * time.Millisecond):
		t.Errorf("Take() didn't return within expected time")
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Create a bucket with capacity 5 and 20 requests per second
	lb := NewLeakyBucket(5, 20)

	// Launch 10 goroutines that all try to access the bucket
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if lb.Allow() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Only 5 should succeed (bucket capacity)
	if successCount != 5 {
		t.Errorf("Expected 5 successful requests, got %d", successCount)
	}
}

func TestEmptyBucket(t *testing.T) {
	// Test with empty bucket (capacity 0)
	lb := NewLeakyBucket(0, 10)

	if lb.Allow() {
		t.Errorf("Bucket with 0 capacity should reject all requests")
	}
}

func TestZeroRateLimit(t *testing.T) {
	// Test with rate limit of 0 (no leaking)
	lb := NewLeakyBucket(5, 0)

	// First 5 should succeed
	for i := 0; i < 5; i++ {
		if !lb.Allow() {
			t.Errorf("Request %d should be allowed", i)
		}
	}

	// 6th should fail
	if lb.Allow() {
		t.Errorf("Request should be rejected after bucket is full")
	}

	// Wait some time, but since rate is 0, no tokens should leak
	lb.lastLeakTime = lb.lastLeakTime.Add(-1 * time.Second)

	// Should still be rejected
	if lb.Allow() {
		t.Errorf("Request should be rejected as no tokens should leak with rate 0")
	}
}

func TestLeakMoreThanTokens(t *testing.T) {
	lb := NewLeakyBucket(5, 10)

	// Add 3 tokens
	for i := 0; i < 3; i++ {
		lb.Allow()
	}

	// Set last leak time to be long ago (should leak all tokens)
	lb.lastLeakTime = lb.lastLeakTime.Add(-1 * time.Second)

	// Internal test of leak function
	lb.mutex.Lock()
	lb.leak()
	tokensAfterLeak := lb.tokens
	lb.mutex.Unlock()

	if tokensAfterLeak != 0 {
		t.Errorf("Expected all tokens to leak, but got %d remaining", tokensAfterLeak)
	}
}

func BenchmarkAllow(b *testing.B) {
	lb := NewLeakyBucket(1000, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb.Allow()
	}
}

func BenchmarkTake(b *testing.B) {
	lb := NewLeakyBucket(b.N, 1000000) // Set high capacity and rate for benchmark

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb.Take()
	}
}
