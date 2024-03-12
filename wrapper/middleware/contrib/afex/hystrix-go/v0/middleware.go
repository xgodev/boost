package hystrix

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/xgodev/boost/errors"
	"github.com/xgodev/boost/log"
	"github.com/xgodev/boost/wrapper"
)

type anyErrorMiddleware[R any] struct {
	name string
}

func (c *anyErrorMiddleware[R]) Exec(ctx *wrapper.AnyErrorContext[R], exec wrapper.AnyErrorExecFunc[R], returnFunc wrapper.AnyErrorReturnFunc[R]) (r R, err error) {
	if err = hystrix.DoC(ctx.GetContext(), c.name,
		func(ctxx context.Context) error {
			r, err = ctx.Next(exec, returnFunc)
			if err != nil {
				return err
			}
			return nil
		},
		func(ctxx context.Context, errr error) error {
			r, err = returnFunc(ctxx, r, errr)
			return err
		}); err != nil {
		return r, errors.Annotate(err, "error during execute hystrix circuit breaker")
	}

	return r, err
}

func NewAnyErrorMiddlewareWithConfig[R any](ctx context.Context, name string, cfg hystrix.CommandConfig) wrapper.AnyErrorMiddleware[R] {
	hystrix.ConfigureCommand(name, cfg)
	hystrix.SetLogger(log.GetLogger())
	return &anyErrorMiddleware[R]{name: name}
}

func NewAnyErrorMiddleware[R any](ctx context.Context, name string) wrapper.AnyErrorMiddleware[R] {
	hystrix.SetLogger(log.GetLogger())
	return &anyErrorMiddleware[R]{name: name}
}

type anyMiddleware[R any] struct {
	name string
}

func (c *anyMiddleware[R]) Exec(ctx *wrapper.AnyContext[R], exec wrapper.AnyExecFunc[R], returnFunc wrapper.AnyReturnFunc[R]) (r R) {
	hystrix.DoC(ctx.GetContext(), c.name,
		func(ctxx context.Context) error {
			r = ctx.Next(exec, returnFunc)
			return nil
		},
		func(ctxx context.Context, errr error) error {
			r = returnFunc(ctxx, r)
			return errr
		})
	return r
}

func NewAnyMiddlewareWithConfig[R any](ctx context.Context, name string, cfg hystrix.CommandConfig) wrapper.AnyMiddleware[R] {
	hystrix.ConfigureCommand(name, cfg)
	hystrix.SetLogger(log.GetLogger())
	return &anyMiddleware[R]{name: name}
}

func NewAnyMiddleware[R any](ctx context.Context, name string) wrapper.AnyMiddleware[R] {
	hystrix.SetLogger(log.GetLogger())
	return &anyMiddleware[R]{name: name}
}

type errorMiddleware struct {
	name string
}

func (c *errorMiddleware) Exec(ctx *wrapper.ErrorContext, exec wrapper.ErrorExecFunc, returnFunc wrapper.ErrorReturnFunc) (err error) {
	err = hystrix.DoC(ctx.GetContext(), c.name,
		func(ctxx context.Context) error {
			return ctx.Next(exec, returnFunc)
		},
		func(ctxx context.Context, errr error) error {
			return returnFunc(ctxx, errr)
		})
	return err
}

func NewErrorMiddlewareWithConfig(ctx context.Context, name string, cfg hystrix.CommandConfig) wrapper.ErrorMiddleware {
	hystrix.ConfigureCommand(name, cfg)
	hystrix.SetLogger(log.GetLogger())
	return &errorMiddleware{name: name}
}

func NewErrorMiddleware(ctx context.Context, name string) wrapper.ErrorMiddleware {
	hystrix.SetLogger(log.GetLogger())
	return &errorMiddleware{name: name}
}
