package timeout

import (
	"context"
	"testing"
	"time"

	"github.com/b0tShaman/go-decorators/api"
)

func MockSlowAPI(delay time.Duration) api.APIFunc {
	return func(ctx context.Context) error {
		select {
		case <-time.After(delay):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func TestTimeout_FailsSlowRequest(t *testing.T) {
	// Timeout: 50ms
	mw := WithTimeout(50 * time.Millisecond)

	// API takes 100ms
	slowAPI := MockSlowAPI(100 * time.Millisecond)
	decorated := mw(slowAPI)

	err := decorated(context.Background())

	if err == nil {
		t.Error("Expected timeout error, got success")
	}
}
