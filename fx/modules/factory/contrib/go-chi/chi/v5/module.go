package chi

import (
	"context"
	"sync"

	c "github.com/go-chi/chi/v5"
	server "github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"github.com/xgodev/boost/fx/modules/extra/multiserver"
	"go.uber.org/fx"
)

type params struct {
	fx.In
	Plugins []chi.Plugin `optional:"true"`
}

var once sync.Once

// Module fx module for the chi v5 app server. It provides *chi.Server and
// *c.Mux and joins the multiserver group, which serves the server on the app
// start and shuts it down on the app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p params) (*chi.Server, error) {
					return chi.NewServer(ctx, p.Plugins...)
				},
				func(srv *chi.Server) *c.Mux {
					return srv.Mux()
				},
			),
			fx.Provide(
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
