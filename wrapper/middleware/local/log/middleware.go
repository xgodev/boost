package log

import (
	"context"
	"github.com/xgodev/boost/log"
	"github.com/xgodev/boost/wrapper"
)

type anyErrorMiddleware[R any] struct{}

func (c *anyErrorMiddleware[R]) Exec(ctx *wrapper.AnyErrorContext[R], exec wrapper.AnyErrorExecFunc[R], returnFunc wrapper.AnyErrorReturnFunc[R]) (r R, err error) {
	l := log.FromContext(ctx.GetContext())
	l.Tracef("executing wrapper %s", ctx.GetName())
	defer l.Debugf("wrapper %s executed", ctx.GetName())
	return ctx.Next(exec, returnFunc)
}

func NewAnyErrorMiddleware[R any](ctx context.Context) wrapper.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{}
}

type anyMiddleware[R any] struct{}

func (c *anyMiddleware[R]) Exec(ctx *wrapper.AnyContext[R], exec wrapper.AnyExecFunc[R], returnFunc wrapper.AnyReturnFunc[R]) (r R) {
	l := log.FromContext(ctx.GetContext())
	l.Tracef("executing wrapper %s", ctx.GetName())
	defer l.Debugf("wrapper %s executed", ctx.GetName())
	return ctx.Next(exec, returnFunc)
}

func NewAnyMiddleware[R any](ctx context.Context) wrapper.AnyMiddleware[R] {
	return &anyMiddleware[R]{}
}

type errorMiddleware struct{}

func (c *errorMiddleware) Exec(ctx *wrapper.ErrorContext, exec wrapper.ErrorExecFunc, returnFunc wrapper.ErrorReturnFunc) (err error) {
	l := log.FromContext(ctx.GetContext())
	l.Tracef("executing wrapper %s", ctx.GetName())
	defer l.Debugf("wrapper %s executed", ctx.GetName())
	return ctx.Next(exec, returnFunc)
}

func NewErrorMiddleware(ctx context.Context) wrapper.ErrorMiddleware {
	return &errorMiddleware{}
}
