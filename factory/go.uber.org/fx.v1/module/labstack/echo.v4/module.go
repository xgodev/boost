package echo

import (
	"context"
	"sync"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/go.uber.org/fx.v1/module/americanas-go/multiserver.v1"
	contextfx "github.com/xgodev/boost/factory/go.uber.org/fx.v1/module/context"
	"github.com/xgodev/boost/factory/labstack/echo.v4"
	server "github.com/xgodev/boost/multiserver"
	"go.uber.org/fx"
)

type params struct {
	fx.In
	Plugins []echo.Plugin `optional:"true"`
}

var once sync.Once

// Module fx module for echo app server.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p params) *echo.Server {
					return echo.NewServer(ctx, p.Plugins...)
				},
				func(srv *echo.Server) *e.Echo {
					return srv.Instance()
				},
			),
			fx.Provide(
				fx.Annotated{
					Group: multiserver.ServersGroupKey,
					Target: func(srv *echo.Server) server.Server {
						return srv
					},
				},
			),
		)
	})

	return options
}
