package ratelimiter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/b0tShaman/go-decorators/api"
)

func MockSuccessAPI(ctx context.Context, req api.Request) api.Response {
	return api.Response{Error: nil}
}

func MockFailAPI(ctx context.Context, req api.Request) api.Response {
	return api.Response{Error: errors.New("simulated failure")}
}

func TestRateLimiter_Blocking(t *testing.T) {
	// Limit: 2 requests, Refill: very slow
	limiter := WithRateLimiting(2, 0.0001)
	decorated := limiter(MockSuccessAPI)

	req := api.Request{UniqueID: "user1"}

	// 1. First request (Should Pass)
	if resp := decorated(context.Background(), req); resp.Error != nil {
		t.Errorf("Expected success, got %v", resp.Error)
	}

	// 2. Second request (Should Pass - bucket had 2)
	if resp := decorated(context.Background(), req); resp.Error != nil {
		t.Errorf("Expected success, got %v", resp.Error)
	}

	// 3. Third request (Should Fail)
	if resp := decorated(context.Background(), req); resp.Error != ErrorRateLimitExceeded {
		t.Errorf("Expected RateLimit error, got %v", resp.Error)
	}
}

func TestRateLimiter_Refill(t *testing.T) {
	// Limit: 1, Refill: 10 per second (0.1s to refill)
	limiter := WithRateLimiting(1, 10.0)
	decorated := limiter(MockSuccessAPI)
	req := api.Request{UniqueID: "userRefill"}

	// Consume token
	decorated(context.Background(), req)

	// Verify blocked
	if resp := decorated(context.Background(), req); resp.Error == nil {
		t.Error("Should have been blocked")
	}

	// Wait for refill
	time.Sleep(150 * time.Millisecond)

	// Should pass now
	if resp := decorated(context.Background(), req); resp.Error != nil {
		t.Errorf("Refill failed: %v", resp.Error)
	}
}
