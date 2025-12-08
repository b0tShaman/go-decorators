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
func WithRateLimiting(limit int, refillRatePerSec float64) Decorator {
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

	uniqueBuckets := make(map[string]*Bucket)
	mu := sync.Mutex{}

	return func(fn APIFunc) APIFunc {
		return func(ctx context.Context, req Request) Response {
			mu.Lock()
			bucket, exists := uniqueBuckets[req.UniqueID]
			if !exists {
				bucket = &Bucket{
					tokens:         float64(limit) - 1,
					lastRefillTime: time.Now(),
				}
				uniqueBuckets[req.UniqueID] = bucket
				mu.Unlock()
				return fn(ctx, req)
			}

			now := time.Now()
			elapsed := now.Sub(bucket.lastRefillTime).Seconds()
			refilledTokens := elapsed * refillRatePerSec
			bucket.tokens = min(float64(limit), bucket.tokens+refilledTokens)
			bucket.lastRefillTime = now

			if bucket.tokens >= 1 {
				bucket.tokens -= 1
				mu.Unlock()
				return fn(ctx, req)
			}
			mu.Unlock()
			return Response{Error: ErrorRateLimitExceeded}
		}
	}
}
