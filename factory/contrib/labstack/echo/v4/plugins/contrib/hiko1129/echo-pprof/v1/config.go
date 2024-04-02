package pprof

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
)

const (
	root    = echo.PluginsRoot + ".pprof"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable pprof integration")
}
