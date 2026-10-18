package server

import (
	"context"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	ms "github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"github.com/xgodev/boost/fx/modules/extra/multiserver"
	"go.uber.org/fx"
)

type params struct {
	fx.In
	Plugins []server.Plugin `optional:"true"`
}

var once sync.Once

// Module fx module for the MCP server. It provides *server.Server and
// *mcp.Server (to register tools) and joins the multiserver group.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p params) (*server.Server, error) {
					return server.NewServer(ctx, p.Plugins...)
				},
				func(srv *server.Server) *mcp.Server {
					return srv.Instance()
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
