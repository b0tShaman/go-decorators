# Go Handy Decorators

This repository provides a collection of **handy decorators** for Go applications. It includes ready-to-use implementations for **Circuit Breaking**, **Rate Limiting**, **Retries**, **Timeouts**, and **Logging**.

## 📦 Usage

**Import into your project:**
```go
import (
    "github.com/b0tShaman/go-decorators/circuitbreaker"
    "github.com/b0tShaman/go-decorators/ratelimiter"
    // ... other packages
)

fn := api.Decorate(API,
    logging.WithLogging(),                                // 1. Outermost: Log entry/exit of every request
    ratelimiter.WithRateLimiting(5, 0.1),                 // 2. Reject excess load cheap and fast
    circuitbreaker.WithCircuitBreaker(5, 10*time.Second), // 3. Protect system if down
    retry.WithRetry(3, 100*time.Millisecond),             // 4. Retry transient errors on healthy system
    timeout.WithTimeout(500*time.Millisecond),            // 5. Per-attempt timeout
)

func API(ctx context.Context, req api.Request) api.Response {
	select {
	case <-time.After(100 * time.Millisecond):
		return api.Response{Error: nil}
	case <-ctx.Done():
		return api.Response{Error: ctx.Err()}
	}
}
