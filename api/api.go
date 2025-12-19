package api

import "context"

type APIFunc func(context.Context) error
type Decorator func(APIFunc) APIFunc

func Decorate(fn APIFunc, decorators ...Decorator) APIFunc {
	for i := len(decorators) - 1; i >= 0; i-- {
		fn = decorators[i](fn)
	}
	return fn
}
