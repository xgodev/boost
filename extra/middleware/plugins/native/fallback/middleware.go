package fallback

import (
	"github.com/xgodev/boost/middleware"
)

type anyErrorMiddleware[R any] struct {
}

func (c *anyErrorMiddleware[R]) Exec(ctx *middleware.AnyErrorContext[R], exec middleware.AnyErrorExecFunc[R], returnFunc middleware.AnyErrorReturnFunc[R]) (R, error) {
	r, err := ctx.Next(exec, returnFunc)
	return returnFunc(ctx.GetContext(), r, err)
}

func NewAnyErrorMiddleware[R any]() middleware.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{}
}

type anyMiddleware[R any] struct {
}

func (c *anyMiddleware[R]) Exec(ctx *middleware.AnyContext[R], exec middleware.AnyExecFunc[R], returnFunc middleware.AnyReturnFunc[R]) R {
	r := ctx.Next(exec, returnFunc)
	return returnFunc(ctx.GetContext(), r)
}

func NewAnyMiddleware[R any]() middleware.AnyMiddleware[R] {
	return &anyMiddleware[R]{}
}

type errorMiddleware struct {
}

func (c *errorMiddleware) Exec(ctx *middleware.ErrorContext, exec middleware.ErrorExecFunc, returnFunc middleware.ErrorReturnFunc) error {
	r := ctx.Next(exec, returnFunc)
	return returnFunc(ctx.GetContext(), r)
}

func NewErrorMiddleware() middleware.ErrorMiddleware {
	return &errorMiddleware{}
}
