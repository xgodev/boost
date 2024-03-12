package fallback

import (
	"github.com/xgodev/boost/grapper"
)

type anyErrorMiddleware[R any] struct {
}

func (c *anyErrorMiddleware[R]) Exec(ctx *grapper.AnyErrorContext[R], exec grapper.AnyErrorExecFunc[R], returnFunc grapper.AnyErrorReturnFunc[R]) (R, error) {
	r, err := ctx.Next(exec, returnFunc)
	return returnFunc(ctx.GetContext(), r, err)
}

func NewAnyErrorMiddleware[R any]() grapper.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{}
}

type anyMiddleware[R any] struct {
}

func (c *anyMiddleware[R]) Exec(ctx *grapper.AnyContext[R], exec grapper.AnyExecFunc[R], returnFunc grapper.AnyReturnFunc[R]) R {
	r := ctx.Next(exec, returnFunc)
	return returnFunc(ctx.GetContext(), r)
}

func NewAnyMiddleware[R any]() grapper.AnyMiddleware[R] {
	return &anyMiddleware[R]{}
}

type errorMiddleware struct {
}

func (c *errorMiddleware) Exec(ctx *grapper.ErrorContext, exec grapper.ErrorExecFunc, returnFunc grapper.ErrorReturnFunc) error {
	r := ctx.Next(exec, returnFunc)
	return returnFunc(ctx.GetContext(), r)
}

func NewErrorMiddleware() grapper.ErrorMiddleware {
	return &errorMiddleware{}
}
