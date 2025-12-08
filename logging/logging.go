package logging

import (
	"context"
	"log"
	"time"

	. "github.com/b0tShaman/go-decorators/api"
)

func WithLogging() Decorator {
	return func(fn APIFunc) APIFunc {
		return func(ctx context.Context, req Request) Response {
			log.Printf("Invoking API with UniqueID: %s", req.UniqueID)
			t := time.Now()
			funcResp := fn(ctx, req)
			log.Printf("API with UniqueID: %s completed in %v", req.UniqueID, time.Since(t))
			return funcResp
		}
	}
}
