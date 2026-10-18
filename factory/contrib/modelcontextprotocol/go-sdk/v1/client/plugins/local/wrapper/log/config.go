package log

import (
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = client.PluginsRoot + ".log"
	enabled = ".enabled"
	level   = ".level"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the plugin keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable logging of each MCP call the client sends")
	config.Add(path+level, "DEBUG", "level of a successful call: INFO/DEBUG/TRACE (errors always log at ERROR)")
}
