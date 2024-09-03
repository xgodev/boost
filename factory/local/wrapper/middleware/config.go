package middleware

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = "boost.factory.local.wrapper.middleware"
	PluginsRoot = root + ".plugins"
	name        = ".name"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+name, "default", "defines default wrapper name")
}
