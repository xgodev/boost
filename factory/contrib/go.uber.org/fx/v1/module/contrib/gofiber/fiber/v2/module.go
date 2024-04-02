package fiber

import (
	"context"
	contextfx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/local/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"sync"

	f "github.com/gofiber/fiber/v2"
	server "github.com/xgodev/boost/extra/multiserver"
	"go.uber.org/fx"
)

type params struct {
	fx.In
	Plugins []fiber.Plugin `optional:"true"`
}

var once sync.Once

// Module fx module for fiber server.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p params) *fiber.Server {
					return fiber.NewServer(ctx, p.Plugins...)
				},
				func(srv *fiber.Server) *f.App {
					return srv.App()
				},
			),
			fx.Provide(
				fx.Annotated{
					Group: multiserver.ServersGroupKey,
					Target: func(srv *fiber.Server) server.Server {
						return srv
					},
				},
			),
		)
	})

	return options
}
