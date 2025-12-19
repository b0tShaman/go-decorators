package ratelimiter

import (
	"context"
	"errors"
	"testing"
	"time"
)

func MockSuccessAPI(ctx context.Context) error {
	return nil
}

func MockFailAPI(ctx context.Context) error {
	return errors.New("simulated failure")
}

func TestRateLimiter_Blocking(t *testing.T) {
	// Limit: 2 requests, Refill: very slow
	limiter := WithRateLimiting(2, 0.0001)
	decorated := limiter(MockSuccessAPI)

	// 1. First request (Should Pass)
	if err := decorated(context.Background()); err != nil {
		t.Errorf("Expected success, got %v", err)
	}

	// 2. Second request (Should Pass - bucket had 2)
	if err := decorated(context.Background()); err != nil {
		t.Errorf("Expected success, got %v", err)
	}

	// 3. Third request (Should Fail)
	if err := decorated(context.Background()); err != ErrorRateLimitExceeded {
		t.Errorf("Expected RateLimit error, got %v", err)
	}
}

func TestRateLimiter_Refill(t *testing.T) {
	// Limit: 1, Refill: 10 per second (0.1s to refill)
	limiter := WithRateLimiting(1, 10.0)
	decorated := limiter(MockSuccessAPI)

	// Consume token
	decorated(context.Background())

	// Verify blocked
	if err := decorated(context.Background()); err == nil {
		t.Error("Should have been blocked")
	}

	// Wait for refill
	time.Sleep(150 * time.Millisecond)

	// Should pass now
	if err := decorated(context.Background()); err != nil {
		t.Errorf("Refill failed: %v", err)
	}
}
