package retry

import (
	"context"
	"fmt"
	"time"

	. "github.com/b0tShaman/go-decorators/api"
)

func WithRetry(times int, delay time.Duration) Decorator {
	return func(fn APIFunc) APIFunc {
		return func(ctx context.Context) error {
			var err error
			for i := 0; i < times; i++ {
				err = fn(ctx)
				if err == nil {
					return err
				}
				time.Sleep(delay)
			}
			return fmt.Errorf("failed after %d attempts: %w", times, err)
		}
	}
}
