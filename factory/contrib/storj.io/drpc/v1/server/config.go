package server

import "github.com/xgodev/boost/wrapper/config"

const (
	root         = "boost.factory.drpc.server"
	port         = ".port"
	collectStats = ".collectStats"
	// PluginsRoot is the config root every drpc server plugin key lives under.
	PluginsRoot = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the drpc server config keys under path.
func ConfigAdd(path string) {
	config.Add(path+port, 9091, "server drpc port")
	config.Add(path+collectStats, false, "collect per rpc stats")
}
