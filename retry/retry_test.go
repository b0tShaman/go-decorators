package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/b0tShaman/go-decorators/api"
)

func TestRetry_SucceedsEventually(t *testing.T) {
	// Fail twice, then succeed
	calls := 0
	flakyAPI := func(ctx context.Context, req api.Request) api.Response {
		calls++
		if calls < 3 {
			return api.Response{Error: errors.New("fail")}
		}
		return api.Response{Error: nil}
	}

	// Retry 3 times
	mw := WithRetry(3, 10*time.Millisecond)
	decorated := mw(flakyAPI)

	resp := decorated(context.Background(), api.Request{})

	if resp.Error != nil {
		t.Errorf("Expected eventual success, got %v", resp.Error)
	}
	if calls != 3 {
		t.Errorf("Expected 3 calls, got %d", calls)
	}
}
