package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_TripsAndRecovers(t *testing.T) {
	// Threshold: 2 failures, Trip Duration: 200ms
	cb := WithCircuitBreaker(2, 200*time.Millisecond)

	// Create a toggleable API: Fail 2 times, then Succeed
	attempts := 0
	flakyAPI := func(ctx context.Context) error {
		attempts++
		if attempts <= 2 {
			return errors.New("oops")
		}
		return nil
	}

	decorated := cb(flakyAPI)

	// 1. Fail twice to trip the breaker
	decorated(context.Background()) // Count = 1
	decorated(context.Background()) // Count = 2 -> TRIPPED

	// 2. Immediate next request should fail FAST (Circuit Open)
	err := decorated(context.Background())
	if err != ErrorCircuitOpen {
		t.Errorf("Expected CircuitOpen error, got %v", err)
	}

	// 3. Wait for timeout (Half-Open Probe)
	time.Sleep(250 * time.Millisecond)

	// 4. Probe request (Should Succeed because flakyAPI is fixed now)
	err = decorated(context.Background())
	if err != nil {
		t.Errorf("Probe should have succeeded, got %v", err)
	}

	// 5. Circuit should be closed now.
	err = decorated(context.Background())
	if err != nil {
		t.Errorf("Circuit should be closed, got %v", err)
	}
}
