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
```

## 🔌 Advanced: Adapting Custom Functions

Use this adapter pattern to wrap custom functions to match the library's APIFunc signature.
```go
// Custom API function to ping a URL
func Ping(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

// Wrapped Ping function which uses library decorators
func PingWrap(url string) (int, error) {
	r := api.Request{UniqueID: url}
	ctx := context.Background()

	type result struct {
		val int
		err error
	}

	var finalResult result

	PingForDec := func(ctx context.Context, r api.Request) api.Response {
		done := make(chan result, 1)
		go func() {
			val, err := Ping(url)
			done <- result{val, err}
		}()

		select {
		case res := <-done:
			finalResult = res
			return api.Response{Error: res.err}
		case <-ctx.Done():
			return api.Response{Error: ctx.Err()}
		}
	}

	pingDecorator := api.Decorate(PingForDec,
		retry.WithRetry(3, 2*time.Second),
		timeout.WithTimeout(1*time.Second),
		logging.WithLogging(),
	)

	res := pingDecorator(ctx, r)
	if res.Error != nil {
		return 0, res.Error
	}
	return finalResult.val, finalResult.err
}
