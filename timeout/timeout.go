package timeout

import (
	"context"
	"errors"
	"time"

	. "github.com/b0tShaman/go-decorators/api"
)

var ErrorTimeout = errors.New("function execution timed out")

// Timeout decorator
func WithTimeout(duration time.Duration) Decorator {
	return func(fn APIFunc) APIFunc {
		return func(ctx context.Context) error {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, duration)
			defer cancel()
			return fn(ctxWithTimeout)
		}
	}
}

// func WithTimeout(duration time.Duration) func(APIFunc) APIFunc {
// 	return func(fn APIFunc) APIFunc {
// 		return func(req Request) Response {
// 			ch := make(chan Response, 1)
// 			go func() {
// 				ch <- fn(req)
// 			}()

// 			select {
// 			case resp := <-ch:
// 				return resp
// 			case <-time.After(duration):
// 				return Response{Error: ErrorTimeout}
// 			}
// 		}
// 	}
// }
