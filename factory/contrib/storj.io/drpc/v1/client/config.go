package client

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = "boost.factory.drpc.client"
	host    = ".host"
	port    = ".port"
	timeout = ".timeout"
	// PluginsRoot is the config root every drpc client plugin key lives under.
	PluginsRoot = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the drpc client config keys under path.
func ConfigAdd(path string) {
	config.Add(path+host, "localhost", "defines host")
	config.Add(path+port, 9091, "defines port")
	config.Add(path+timeout, 10*time.Second, "defines dial timeout")
}
