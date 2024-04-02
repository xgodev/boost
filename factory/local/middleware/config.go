package middleware

import (
	"github.com/xgodev/boost/config"
)

const (
	root        = "boost.factory.wrapper"
	PluginsRoot = root + ".plugins"
	name        = ".name"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+name, "default", "defines default wrapper name")
}
