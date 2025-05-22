package token_bucket

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens         float64    // Current number of tokens
	maxTokens      float64    // Maximum tokens allowed
	refillRate     float64    // Tokens added per second
	lastRefillTime time.Time  // Last time tokens were refilled
	mutex          sync.Mutex // For thread safety
}

func NewTokenBucket(maxTokens, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	duration := now.Sub(tb.lastRefillTime).Seconds()

	// Calculate tokens to add based on time elapsed
	tokensToAdd := duration * tb.refillRate

	// Update token count, but don't exceed max capacity
	tb.tokens = tb.tokens + tokensToAdd
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}

	// Update last refill time
	tb.lastRefillTime = now
}

func (tb *TokenBucket) Allow(tokensRequired float64) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	// Refill tokens based on elapsed time
	tb.refill()

	// Check if we have enough tokens
	if tb.tokens >= tokensRequired {
		tb.tokens -= tokensRequired
		return true
	}

	return false
}
