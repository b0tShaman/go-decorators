package api

import "context"

type APIFunc func(context.Context, Request) Response
type Decorator func(APIFunc) APIFunc

type Request struct {
	UniqueID string
}

type Response struct {
	Error error
}

func Decorate(fn APIFunc, decorators ...Decorator) APIFunc {
	for i := len(decorators) - 1; i >= 0; i-- {
		fn = decorators[i](fn)
	}
	return fn
}
