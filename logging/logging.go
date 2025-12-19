package logging

import (
	"context"
	"log"
	"time"

	. "github.com/b0tShaman/go-decorators/api"
)

func WithLogging() Decorator {
	return func(fn APIFunc) APIFunc {
		return func(ctx context.Context) error {
			log.Printf("Invoking API")
			t := time.Now()
			err := fn(ctx)
			log.Printf("API completed in %v", time.Since(t))
			return err
		}
	}
}
