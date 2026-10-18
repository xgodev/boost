package client

import (
	"context"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

type params struct {
	fx.In
	Plugins []client.Plugin `optional:"true"`
}

var once sync.Once

// Module fx module for the MCP client. It provides *client.Client and its
// *mcp.ClientSession; the session is closed on app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle, p params) (*client.Client, error) {
					cli, err := client.NewClient(ctx, p.Plugins...)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(cli.Close))
					return cli, nil
				},
				func(cli *client.Client) *mcp.ClientSession {
					return cli.Session()
				},
			),
		)
	})

	return options
}
