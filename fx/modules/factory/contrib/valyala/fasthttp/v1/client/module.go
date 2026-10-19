package client

import (
	"context"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/xgodev/boost/factory/contrib/valyala/fasthttp/v1/client"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

type params struct {
	fx.In
	Plugins []client.Plugin `optional:"true"`
}

var once sync.Once

// Module fx module for the fasthttp client. It provides *fasthttp.Client.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p params) (*fasthttp.Client, error) {
					return client.NewClient(ctx, p.Plugins...)
				},
			),
		)
	})

	return options
}
