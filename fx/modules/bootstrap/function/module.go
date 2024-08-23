package function

import (
	"context"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/extra/middleware"
	fxcontext "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
	"sync"
)

const (
	BSFunctionAdaptersGroupKey    = "boostrap.function.adapters"
	BSFunctionMiddlewaresGroupKey = "boostrap.function.middlewares"
)

type params[T any] struct {
	fx.In
	Adapters    []function.CmdFunc[T]              `group:"boostrap.function.adapters"`
	Middlewares []middleware.AnyErrorMiddleware[T] `group:"boostrap.function.middlewares"`
}

var once sync.Once

func Module[T any]() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			fxcontext.Module(),
			fx.Provide(
				func(ctx context.Context, p params[T]) (*function.Function[T], error) {
					return function.New[T](p.Middlewares...)
				}),
			fx.Invoke(
				func(ctx context.Context, p params[T], hdl function.Handler[T], fn *function.Function[T]) error {
					return fn.Run(ctx, hdl, p.Adapters...)
				},
			),
		)
	})

	return options
}
