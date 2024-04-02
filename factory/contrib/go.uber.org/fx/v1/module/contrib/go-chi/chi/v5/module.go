package chi

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	contextfx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/local/extra/multiserver"
	"sync"

	c "github.com/go-chi/chi/v5"
	server "github.com/xgodev/boost/extra/multiserver"
	"go.uber.org/fx"
)

var once sync.Once

type params struct {
	fx.In
	Plugins []chi.Plugin `optional:"true"`
}

// Module fx module for chi.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p params) *chi.Server {
					return chi.NewServer(ctx, p.Plugins...)
				},
				func(srv *chi.Server) *c.Mux {
					return srv.Mux()
				},
				fx.Annotated{
					Group: multiserver.ServersGroupKey,
					Target: func(srv *chi.Server) server.Server {
						return srv
					},
				},
			),
		)
	})

	return options
}
