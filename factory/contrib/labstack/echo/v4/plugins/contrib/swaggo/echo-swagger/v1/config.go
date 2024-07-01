package swagger

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = echo.PluginsRoot + ".swagger"
	enabled = ".enabled"
	route   = ".route"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable swagger integration")
	config.Add(path+route, "/swagger/*", "define swagger metrics url")
}
