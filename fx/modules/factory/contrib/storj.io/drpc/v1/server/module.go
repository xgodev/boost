package server

import (
	"context"
	"sync"

	ms "github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/storj.io/drpc/v1/server"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"github.com/xgodev/boost/fx/modules/extra/multiserver"
	"go.uber.org/fx"
	"storj.io/drpc/drpcmux"
)

type params struct {
	fx.In
	Plugins []server.Plugin `optional:"true"`
}

var once sync.Once

// Module fx module for the drpc server. It provides *server.Server and
// *drpcmux.Mux and joins the multiserver group, which serves the server on
// the app start and shuts it down on the app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p params) (*server.Server, error) {
					return server.NewServer(ctx, p.Plugins...)
				},
				func(srv *server.Server) *drpcmux.Mux {
					return srv.Mux()
				},
			),
			fx.Provide(
				fx.Annotated{
					Group: multiserver.ServersGroupKey,
					Target: func(srv *server.Server) ms.Server {
						return srv
					},
				},
			),
		)
	})

	return options
}
