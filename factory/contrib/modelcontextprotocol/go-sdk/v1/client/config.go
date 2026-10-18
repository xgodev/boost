package client

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

const (
	root       = "boost.factory.mcp.client"
	name       = ".name"
	version    = ".version"
	transport  = ".transport"
	endpoint   = ".endpoint"
	maxRetries = ".maxRetries"
	command    = ".command"
	args       = ".args"
	keepAlive  = ".keepAlive"

	// PluginsRoot is where each client plugin package roots its keys.
	PluginsRoot = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the client keys under path, one path per MCP server
// the application calls.
func ConfigAdd(path string) {
	config.Add(path+name, "", "implementation name sent on initialize; empty uses boost.applicationName")
	config.Add(path+version, "", "implementation version sent on initialize; empty uses the binary's main module version (0.0.0 when unstamped)")
	config.Add(path+transport, TransportHTTP, "transport: HTTP (streamable HTTP endpoint) or COMMAND (subprocess over stdio)")
	config.Add(path+endpoint, "http://localhost:8090/mcp", "HTTP transport server URL")
	config.Add(path+maxRetries, 5, "HTTP transport reconnect attempts; negative disables")
	config.Add(path+command, "", "COMMAND transport executable")
	config.Add(path+args, []string{}, "COMMAND transport arguments")
	config.Add(path+keepAlive, time.Duration(0), "ping interval that closes a dead session; 0 disables")
}
