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

// Core Type Definitions (in package api)
type Request struct {
    UniqueID string
}

type Response struct {
    Error error
}

type APIFunc func(context.Context, Request) Response
type Decorator func(APIFunc) APIFunc