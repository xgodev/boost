package fallback

import (
	"github.com/xgodev/boost/wrapper"
)

type anyErrorMiddleware[R any] struct {
}

func (c *anyErrorMiddleware[R]) Exec(ctx *wrapper.AnyErrorContext[R], exec wrapper.AnyErrorExecFunc[R], returnFunc wrapper.AnyErrorReturnFunc[R]) (R, error) {
	r, err := ctx.Next(exec, returnFunc)
	return returnFunc(ctx.GetContext(), r, err)
}

func NewAnyErrorMiddleware[R any]() wrapper.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{}
}

type anyMiddleware[R any] struct {
}

func (c *anyMiddleware[R]) Exec(ctx *wrapper.AnyContext[R], exec wrapper.AnyExecFunc[R], returnFunc wrapper.AnyReturnFunc[R]) R {
	r := ctx.Next(exec, returnFunc)
	return returnFunc(ctx.GetContext(), r)
}

func NewAnyMiddleware[R any]() wrapper.AnyMiddleware[R] {
	return &anyMiddleware[R]{}
}

type errorMiddleware struct {
}

func (c *errorMiddleware) Exec(ctx *wrapper.ErrorContext, exec wrapper.ErrorExecFunc, returnFunc wrapper.ErrorReturnFunc) error {
	r := ctx.Next(exec, returnFunc)
	return returnFunc(ctx.GetContext(), r)
}

func NewErrorMiddleware() wrapper.ErrorMiddleware {
	return &errorMiddleware{}
}
