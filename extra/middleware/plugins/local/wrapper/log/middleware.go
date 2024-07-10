package log

import (
	"context"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/wrapper/log"
)

type anyErrorMiddleware[R any] struct{}

func (c *anyErrorMiddleware[R]) Exec(ctx *middleware.AnyErrorContext[R], exec middleware.AnyErrorExecFunc[R], returnFunc middleware.AnyErrorReturnFunc[R]) (r R, err error) {
	l := log.FromContext(ctx.GetContext())
	l.Tracef("executing wrapper %s", ctx.GetName())
	defer l.Debugf("wrapper %s executed", ctx.GetName())
	return ctx.Next(exec, returnFunc)
}

func NewAnyErrorMiddleware[R any](ctx context.Context) middleware.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{}
}

type anyMiddleware[R any] struct{}

func (c *anyMiddleware[R]) Exec(ctx *middleware.AnyContext[R], exec middleware.AnyExecFunc[R], returnFunc middleware.AnyReturnFunc[R]) (r R) {
	l := log.FromContext(ctx.GetContext())
	l.Tracef("executing wrapper %s", ctx.GetName())
	defer l.Debugf("wrapper %s executed", ctx.GetName())
	return ctx.Next(exec, returnFunc)
}

func NewAnyMiddleware[R any](ctx context.Context) middleware.AnyMiddleware[R] {
	return &anyMiddleware[R]{}
}

type errorMiddleware struct{}

func (c *errorMiddleware) Exec(ctx *middleware.ErrorContext, exec middleware.ErrorExecFunc, returnFunc middleware.ErrorReturnFunc) (err error) {
	l := log.FromContext(ctx.GetContext())
	l.Tracef("executing wrapper %s", ctx.GetName())
	defer l.Debugf("wrapper %s executed", ctx.GetName())
	return ctx.Next(exec, returnFunc)
}

func NewErrorMiddleware(ctx context.Context) middleware.ErrorMiddleware {
	return &errorMiddleware{}
}
