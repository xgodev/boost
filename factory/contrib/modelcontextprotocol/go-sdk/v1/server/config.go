package server

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

const (
	root           = "boost.factory.mcp.server"
	name           = ".name"
	version        = ".version"
	instructions   = ".instructions"
	transport      = ".transport"
	port           = ".port"
	httpPath       = ".path"
	stateless      = ".stateless"
	jsonResponse   = ".jsonResponse"
	sessionTimeout = ".sessionTimeout"
	keepAlive      = ".keepAlive"

	// PluginsRoot is where each server plugin package roots its keys.
	PluginsRoot = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the server keys under path, so one binary can run a
// second MCP server (another port, another tool set) from its own path.
func ConfigAdd(path string) {
	config.Add(path+name, "", "implementation name sent on initialize; empty uses boost.applicationName")
	config.Add(path+version, "", "implementation version sent on initialize; empty uses the binary's main module version (0.0.0 when unstamped)")
	config.Add(path+instructions, "", "instructions the client may add to its system prompt")
	config.Add(path+transport, TransportHTTP, "transport: HTTP (streamable HTTP on port/path) or STDIO")
	// TODO(maintainer-review): 8090 keeps clear of echo (8080) and gRPC (9090)
	// so the three can share a multiserver without extra config.
	config.Add(path+port, 8090, "HTTP transport port")
	config.Add(path+httpPath, "/mcp", "HTTP transport path")
	config.Add(path+stateless, false, "HTTP transport without session state (no server-to-client requests)")
	config.Add(path+jsonResponse, false, "HTTP transport answers application/json instead of text/event-stream")
	config.Add(path+sessionTimeout, time.Duration(0), "HTTP transport idle session timeout; 0 never closes")
	config.Add(path+keepAlive, time.Duration(0), "ping interval that closes dead sessions; 0 disables")
}
