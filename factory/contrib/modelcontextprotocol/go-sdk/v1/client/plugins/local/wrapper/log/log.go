// Package log logs every MCP call the client sends through wrapper/log:
// method, tool, duration and error.
package log

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/internal/middleware"
)

// Register installs the log middleware on cli; a client.Plugin.
func Register(_ context.Context, cli *client.Client) error {
	o, err := NewOptions()
	if err != nil {
		return err
	}
	if !o.Enabled {
		return nil
	}
	cli.Instance().AddSendingMiddleware(middleware.Log(o.Level))
	return nil
}
