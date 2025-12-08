package retry

import (
	"context"
	"time"

	. "github.com/b0tShaman/go-decorators/api"
)

func WithRetry(times int, delay time.Duration) Decorator {
	return func(fn APIFunc) APIFunc {
		return func(ctx context.Context, req Request) Response {
			var resp Response
			for i := 0; i < times; i++ {
				resp = fn(ctx, req)
				if resp.Error == nil {
					return resp
				}
				time.Sleep(delay)
			}
			return resp
		}
	}
}
