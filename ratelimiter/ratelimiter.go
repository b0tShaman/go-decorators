package ratelimiter

import (
	"context"
	"errors"
	"sync"
	"time"

	. "github.com/b0tShaman/go-decorators/api"
)

var ErrorRateLimitExceeded = errors.New("rate limit exceeded")

// Token Bucket rate limiting decorator
func WithRateLimiting(limit float64, refillRatePerSec float64) Decorator {
	type Bucket struct {
		tokens         float64
		lastRefillTime time.Time
	}

	min := func(a, b float64) float64 {
		if a < b {
			return a
		}
		return b
	}

	bucket := &Bucket{
		tokens:         limit,
		lastRefillTime: time.Now(),
	}
	mu := sync.Mutex{}

	return func(fn APIFunc) APIFunc {
		return func(ctx context.Context) error {
			mu.Lock()

			now := time.Now()
			elapsed := now.Sub(bucket.lastRefillTime).Seconds()

			bucket.tokens = min(limit, bucket.tokens+elapsed*refillRatePerSec)
			bucket.lastRefillTime = now

			if bucket.tokens >= 1 {
				bucket.tokens -= 1
				mu.Unlock()
				return fn(ctx)
			}
			mu.Unlock()
			return ErrorRateLimitExceeded
		}
	}
}
