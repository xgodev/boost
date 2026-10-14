package datadog

import (
	"context"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/xgodev/boost/extra/middleware"
)

type anyErrorMiddleware[R any] struct {
	name string
	tp   string
}

func (c *anyErrorMiddleware[R]) Exec(ctx *middleware.AnyErrorContext[R], exec middleware.AnyErrorExecFunc[R], returnFunc middleware.AnyErrorReturnFunc[R]) (r R, err error) {
	span, sctx := tracer.StartSpanFromContext(
		ctx.GetContext(),
		c.name,
		tracer.SpanType(c.tp),
	)
	defer span.Finish()

	ctx.SetContext(sctx)

	return ctx.Next(exec, returnFunc)
}

func NewAnyErrorMiddleware[R any](ctx context.Context, name string, tp string) middleware.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{name: name, tp: tp}
}

type anyMiddleware[R any] struct {
	name string
	tp   string
}

func (c *anyMiddleware[R]) Exec(ctx *middleware.AnyContext[R], exec middleware.AnyExecFunc[R], returnFunc middleware.AnyReturnFunc[R]) (r R) {
	span, sctx := tracer.StartSpanFromContext(
		ctx.GetContext(),
		c.name,
		tracer.SpanType(c.tp),
	)
	defer span.Finish()

	ctx.SetContext(sctx)

	return ctx.Next(exec, returnFunc)
}

func NewAnyMiddleware[R any](ctx context.Context, name string, tp string) middleware.AnyMiddleware[R] {
	return &anyMiddleware[R]{name: name, tp: tp}
}

type errorMiddleware struct {
	name string
	tp   string
}

func (c *errorMiddleware) Exec(ctx *middleware.ErrorContext, exec middleware.ErrorExecFunc, returnFunc middleware.ErrorReturnFunc) (err error) {
	span, sctx := tracer.StartSpanFromContext(
		ctx.GetContext(),
		c.name,
		tracer.SpanType(c.tp),
	)
	defer span.Finish()

	ctx.SetContext(sctx)

	return ctx.Next(exec, returnFunc)
}

func NewErrorMiddleware(ctx context.Context, name string, tp string) middleware.ErrorMiddleware {
	return &errorMiddleware{name: name, tp: tp}
}
