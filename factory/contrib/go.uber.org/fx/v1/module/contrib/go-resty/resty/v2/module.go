package resty

import (
	"context"
	r "github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
	fxcontext "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"go.uber.org/fx"
	"sync"
)

const (
	PluginsGroupKey = "factory.resty.plugins"
)

type params struct {
	fx.In
	Plugins []resty.Plugin `optional:"true"`
}

var once sync.Once

func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Module("resty",
			fxcontext.Module(),
			fx.Provide(
				func(ctx context.Context, p params) (*r.Client, error) {
					return resty.NewClient(ctx, p.Plugins...) //nolint:wrapcheck  // This must be resolved in another task
				},
			),
		)
	})

	return options
}
