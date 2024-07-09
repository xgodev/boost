package multiserver

import (
	"context"
	contextfx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
	"sync"

	c "github.com/spf13/cobra"
	server "github.com/xgodev/boost/extra/multiserver"
	"go.uber.org/fx"
)

const (
	ServersGroupKey = "extra.multiserver.servers"
)

type params struct {
	fx.In
	Servers []server.Server `group:"extra.multiserver.servers"`
}

var once sync.Once

// Module fx module for multiserver.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Invoke(
				func(ctx context.Context, p params) error {

					return cobra.Run(
						&c.Command{
							Run: func(cmd *c.Command, args []string) {
								server.Serve(ctx, p.Servers...)
							},
						},
					)

				},
			),
		)
	})

	return options
}
