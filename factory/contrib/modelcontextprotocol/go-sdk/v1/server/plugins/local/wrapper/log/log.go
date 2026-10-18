// Package log logs every MCP call the server receives through wrapper/log:
// method, tool, duration and error.
package log

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/internal/middleware"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
)

// Register installs the log middleware on srv; a server.Plugin.
func Register(_ context.Context, srv *server.Server) error {
	o, err := NewOptions()
	if err != nil {
		return err
	}
	if !o.Enabled {
		return nil
	}
	srv.Instance().AddReceivingMiddleware(middleware.Log(o.Level))
	return nil
}
