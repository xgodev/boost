package otel

import (
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = client.PluginsRoot + ".otel"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the plugin keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable an OpenTelemetry span for each MCP call the client sends")
}
