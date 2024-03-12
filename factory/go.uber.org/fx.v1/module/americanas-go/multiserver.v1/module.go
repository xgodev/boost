package multiserver

import (
	"context"
	"sync"

	c "github.com/spf13/cobra"
	contextfx "github.com/xgodev/boost/factory/go.uber.org/fx.v1/module/context"
	"github.com/xgodev/boost/factory/spf13/cobra.v1"
	server "github.com/xgodev/boost/multiserver"
	"go.uber.org/fx"
)

const (
	ServersGroupKey = "_gi_server_servers_"
)

type srvParams struct {
	fx.In
	Servers []server.Server `group:"_gi_server_servers_"`
}

var once sync.Once

// Module fx module for multiserver.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Invoke(
				func(ctx context.Context, p srvParams) error {

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
