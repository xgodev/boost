// Package mcp holds what the mcp server of app api needs besides its
// handlers.
package mcp

import (
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server/plugins/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server/plugins/local/wrapper/log"
)

// Plugins log the tool calls and trace them with OpenTelemetry.
// @Provide(index=0)
func Plugins() []server.Plugin {
	return []server.Plugin{
		log.Register,
		otel.Register,
	}
}
