package function

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/bootstrap/function"
	fxcontext "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"github.com/xgodev/boost/middleware"
	"go.uber.org/fx"
	"sync"
)

const (
	BSFunctionAdaptersGroupKey    = "boostrap.function.adapters"
	BSFunctionMiddlewaresGroupKey = "boostrap.function.middlewares"
)

type params struct {
	fx.In
	Adapters    []function.CmdFunc                            `group:"boostrap.function.adapters"`
	Middlewares []middleware.AnyErrorMiddleware[*event.Event] `group:"boostrap.function.middlewares" optional:"true"`
}

var once sync.Once

func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			fxcontext.Module(),
			fx.Provide(
				func(ctx context.Context, p params) *function.Function {
					return function.New(p.Middlewares...)
				}),
			fx.Invoke(
				func(ctx context.Context, p params, hdl function.Handler, fn *function.Function) error {
					return fn.Run(ctx, hdl, p.Adapters...)
				},
			),
		)
	})

	return options
}
