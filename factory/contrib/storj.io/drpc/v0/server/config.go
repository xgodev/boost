package server

import "github.com/xgodev/boost/config"

const (
	root        = "boost.factory.drpc.server"
	port        = ".port"
	PluginsRoot = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+port, 9091, "server drpc port")
}
