package log

import (
	"context"
	"github.com/xgodev/boost/grapper"
	"github.com/xgodev/boost/log"
)

type anyErrorMiddleware[R any] struct{}

func (c *anyErrorMiddleware[R]) Exec(ctx *grapper.AnyErrorContext[R], exec grapper.AnyErrorExecFunc[R], returnFunc grapper.AnyErrorReturnFunc[R]) (r R, err error) {
	l := log.FromContext(ctx.GetContext())
	l.Tracef("executing wrapper %s", ctx.GetName())
	defer l.Debugf("wrapper %s executed", ctx.GetName())
	return ctx.Next(exec, returnFunc)
}

func NewAnyErrorMiddleware[R any](ctx context.Context) grapper.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{}
}

type anyMiddleware[R any] struct{}

func (c *anyMiddleware[R]) Exec(ctx *grapper.AnyContext[R], exec grapper.AnyExecFunc[R], returnFunc grapper.AnyReturnFunc[R]) (r R) {
	l := log.FromContext(ctx.GetContext())
	l.Tracef("executing wrapper %s", ctx.GetName())
	defer l.Debugf("wrapper %s executed", ctx.GetName())
	return ctx.Next(exec, returnFunc)
}

func NewAnyMiddleware[R any](ctx context.Context) grapper.AnyMiddleware[R] {
	return &anyMiddleware[R]{}
}

type errorMiddleware struct{}

func (c *errorMiddleware) Exec(ctx *grapper.ErrorContext, exec grapper.ErrorExecFunc, returnFunc grapper.ErrorReturnFunc) (err error) {
	l := log.FromContext(ctx.GetContext())
	l.Tracef("executing wrapper %s", ctx.GetName())
	defer l.Debugf("wrapper %s executed", ctx.GetName())
	return ctx.Next(exec, returnFunc)
}

func NewErrorMiddleware(ctx context.Context) grapper.ErrorMiddleware {
	return &errorMiddleware{}
}
