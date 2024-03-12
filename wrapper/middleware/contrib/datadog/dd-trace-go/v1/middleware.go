package datadog

import (
	"context"
	"github.com/xgodev/boost/wrapper"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type anyErrorMiddleware[R any] struct {
	name string
	tp   string
}

func (c *anyErrorMiddleware[R]) Exec(ctx *wrapper.AnyErrorContext[R], exec wrapper.AnyErrorExecFunc[R], returnFunc wrapper.AnyErrorReturnFunc[R]) (r R, err error) {
	span, sctx := tracer.StartSpanFromContext(
		ctx.GetContext(),
		c.name,
		tracer.SpanType(c.tp),
	)
	defer span.Finish()

	ctx.SetContext(sctx)

	return ctx.Next(exec, returnFunc)
}

func NewAnyErrorMiddleware[R any](ctx context.Context, name string, tp string) wrapper.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{name: name, tp: tp}
}

type anyMiddleware[R any] struct {
	name string
	tp   string
}

func (c *anyMiddleware[R]) Exec(ctx *wrapper.AnyContext[R], exec wrapper.AnyExecFunc[R], returnFunc wrapper.AnyReturnFunc[R]) (r R) {
	span, sctx := tracer.StartSpanFromContext(
		ctx.GetContext(),
		c.name,
		tracer.SpanType(c.tp),
	)
	defer span.Finish()

	ctx.SetContext(sctx)

	return ctx.Next(exec, returnFunc)
}

func NewAnyMiddleware[R any](ctx context.Context, name string, tp string) wrapper.AnyMiddleware[R] {
	return &anyMiddleware[R]{name: name, tp: tp}
}

type errorMiddleware struct {
	name string
	tp   string
}

func (c *errorMiddleware) Exec(ctx *wrapper.ErrorContext, exec wrapper.ErrorExecFunc, returnFunc wrapper.ErrorReturnFunc) (err error) {
	span, sctx := tracer.StartSpanFromContext(
		ctx.GetContext(),
		c.name,
		tracer.SpanType(c.tp),
	)
	defer span.Finish()

	ctx.SetContext(sctx)

	return ctx.Next(exec, returnFunc)
}

func NewErrorMiddleware(ctx context.Context, name string, tp string) wrapper.ErrorMiddleware {
	return &errorMiddleware{name: name, tp: tp}
}
