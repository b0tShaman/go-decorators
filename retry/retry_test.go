package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetry_SucceedsEventually(t *testing.T) {
	// Fail twice, then succeed
	calls := 0
	flakyAPI := func(ctx context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("fail")
		}
		return nil
	}

	// Retry 3 times
	mw := WithRetry(3, 10*time.Millisecond)
	decorated := mw(flakyAPI)

	err := decorated(context.Background())

	if err != nil {
		t.Errorf("Expected eventual success, got %v", err)
	}
	if calls != 3 {
		t.Errorf("Expected 3 calls, got %d", calls)
	}
}
