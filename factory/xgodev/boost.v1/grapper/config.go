package grapper

import (
	"github.com/xgodev/boost/config"
)

const (
	root        = "ignite.grapper"
	PluginsRoot = root + ".plugins"
	name        = ".name"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+name, "default", "defines default wrapper name")
}
