package client

import (
	"context"
	"sync"

	c "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var once sync.Once

type params struct {
	fx.In
	Plugins []c.Plugin `optional:"true"`
}

// Module fx module for grpc client connection. The connection is closed on app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle, p params) *grpc.ClientConn {
					conn := c.NewClientConn(ctx, p.Plugins...)
					lc.Append(fx.StopHook(conn.Close))
					return conn
				},
			),
		)
	})

	return options
}
