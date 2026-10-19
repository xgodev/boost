// Package mcp holds what the mcp server of app api needs besides its tools.
package mcp

import (
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server/plugins/local/wrapper/log"
)

// Plugins log every tool call of the mcp server.
// @Provide(index=0)
func Plugins() []server.Plugin {
	return []server.Plugin{log.Register}
}
