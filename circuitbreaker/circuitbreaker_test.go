package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/b0tShaman/go-decorators/api"
)

func TestCircuitBreaker_TripsAndRecovers(t *testing.T) {
	// Threshold: 2 failures, Trip Duration: 200ms
	cb := WithCircuitBreaker(2, 200*time.Millisecond)

	// Create a toggleable API: Fail 2 times, then Succeed
	attempts := 0
	flakyAPI := func(ctx context.Context, req api.Request) api.Response {
		attempts++
		if attempts <= 2 {
			return api.Response{Error: errors.New("oops")}
		}
		return api.Response{Error: nil}
	}

	decorated := cb(flakyAPI)
	req := api.Request{UniqueID: "cb_test"}

	// 1. Fail twice to trip the breaker
	decorated(context.Background(), req) // Count = 1
	decorated(context.Background(), req) // Count = 2 -> TRIPPED

	// 2. Immediate next request should fail FAST (Circuit Open)
	resp := decorated(context.Background(), req)
	if resp.Error != ErrorCircuitOpen {
		t.Errorf("Expected CircuitOpen error, got %v", resp.Error)
	}

	// 3. Wait for timeout (Half-Open Probe)
	time.Sleep(250 * time.Millisecond)

	// 4. Probe request (Should Succeed because flakyAPI is fixed now)
	resp = decorated(context.Background(), req)
	if resp.Error != nil {
		t.Errorf("Probe should have succeeded, got %v", resp.Error)
	}

	// 5. Circuit should be closed now.
	resp = decorated(context.Background(), req)
	if resp.Error != nil {
		t.Errorf("Circuit should be closed, got %v", resp.Error)
	}
}
