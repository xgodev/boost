package client

import (
	"context"
	"sync"

	"github.com/xgodev/boost/factory/contrib/storj.io/drpc/v1/client"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
	"storj.io/drpc/drpcconn"
)

type params struct {
	fx.In
	Plugins []client.Plugin `optional:"true"`
}

var once sync.Once

// Module fx module for the drpc client. It provides *drpcconn.Conn, dialed
// on the app start, and closes it on the app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle, p params) (*drpcconn.Conn, error) {
					conn, err := client.NewClientConn(ctx, p.Plugins...)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(conn.Close))
					return conn, nil
				},
			),
		)
	})

	return options
}
