package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"

	. "github.com/b0tShaman/go-decorators/api"
)

var ErrorCircuitOpen = errors.New("circuit breaker is open")
var ErrorCircuitHalfOpen = errors.New("circuit breaker is half-open")


// Circuit Breaker decorator
func WithCircuitBreaker(failureThreshold int, tripDuration time.Duration) Decorator {
	const (
		StateOpen = iota
		StateHalfOpen
		StateClosed
	)

	state := StateClosed
	failureCount := 0
	lastTripped := time.Time{}

	mu := sync.Mutex{}
	return func(fn APIFunc) APIFunc {
		return func(ctx context.Context, req Request) Response {
			var resp Response
			mu.Lock()
			if state == StateOpen {
				if time.Since(lastTripped) <= tripDuration {
					mu.Unlock()
					return Response{Error: ErrorCircuitOpen}
				}
				state = StateHalfOpen
				failureCount = failureThreshold - 1
			} else if state == StateHalfOpen {
				mu.Unlock()
				return Response{Error: ErrorCircuitHalfOpen}
			}

			mu.Unlock()

			resp = fn(ctx, req)
			mu.Lock()
			defer mu.Unlock()

			if resp.Error != nil {
				failureCount++
				if failureCount >= failureThreshold {
					state = StateOpen
					lastTripped = time.Now()
				}
				return resp
			}
			// Success
			failureCount = 0
			state = StateClosed
			return resp
		}
	}
}
