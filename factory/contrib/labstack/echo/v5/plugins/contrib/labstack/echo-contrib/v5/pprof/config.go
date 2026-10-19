package pprof

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = echo.PluginsRoot + ".pprof"
	enabled = ".enabled"
	prefix  = ".prefix"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the pprof plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable pprof routes")
	config.Add(path+prefix, "/debug/pprof", "define pprof routes prefix")
}
