package client

import (
	"github.com/xgodev/boost/config"
)

const (
	root        = "boost.factory.drpc.client"
	PluginsRoot = root + ".plugins"
	host        = ".host"
	port        = ".port"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+host, "localhost", "defines host")
	config.Add(path+port, 9091, "defines port")
}
