package timeout

import (
	"context"
	"testing"
	"time"

	"github.com/b0tShaman/go-decorators/api"
)

func MockSlowAPI(delay time.Duration) api.APIFunc {
	return func(ctx context.Context, req api.Request) api.Response {
		select {
		case <-time.After(delay):
			return api.Response{Error: nil}
		case <-ctx.Done():
			return api.Response{Error: ctx.Err()}
		}
	}
}

func TestTimeout_FailsSlowRequest(t *testing.T) {
	// Timeout: 50ms
	mw := WithTimeout(50 * time.Millisecond)

	// API takes 100ms
	slowAPI := MockSlowAPI(100 * time.Millisecond)
	decorated := mw(slowAPI)

	resp := decorated(context.Background(), api.Request{})

	if resp.Error == nil {
		t.Error("Expected timeout error, got success")
	}
}
