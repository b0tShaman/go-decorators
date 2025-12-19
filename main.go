package main

import (
	"context"
	"net/http"
	"time"

	"github.com/b0tShaman/go-decorators/api"
	"github.com/b0tShaman/go-decorators/circuitbreaker"
	"github.com/b0tShaman/go-decorators/logging"
	"github.com/b0tShaman/go-decorators/ratelimiter"
	"github.com/b0tShaman/go-decorators/retry"
	"github.com/b0tShaman/go-decorators/timeout"
)

func main() {
	mux := http.NewServeMux()

	fn := api.Decorate(API,
		logging.WithLogging(),                                // 1. Outermost: Log entry/exit of every request
		ratelimiter.WithRateLimiting(5, 0.1),                 // 2. Reject excess load cheap and fast
		circuitbreaker.WithCircuitBreaker(5, 10*time.Second), // 3. Protect system if down
		retry.WithRetry(3, 100*time.Millisecond),             // 4. Retry transient errors on healthy system
		timeout.WithTimeout(500*time.Millisecond),            // 5. Per-attempt timeout
	)

	mux.HandleFunc("POST /api", func(w http.ResponseWriter, r *http.Request) {
		err := fn(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", mux)
}

func API(ctx context.Context) error {
	select {
	case <-time.After(100 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
