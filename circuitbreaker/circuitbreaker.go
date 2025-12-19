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
		return func(ctx context.Context) error {
			mu.Lock()
			if state == StateOpen {
				if time.Since(lastTripped) <= tripDuration {
					mu.Unlock()
					return ErrorCircuitOpen
				}
				state = StateHalfOpen
				failureCount = failureThreshold - 1 // Allow one attempt
			} else if state == StateHalfOpen {
				mu.Unlock()
				return ErrorCircuitHalfOpen
			}

			mu.Unlock()

			err := fn(ctx)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				failureCount++
				if failureCount >= failureThreshold {
					state = StateOpen
					lastTripped = time.Now()
				}
				return err
			}
			// Success
			failureCount = 0
			state = StateClosed
			return nil
		}
	}
}
