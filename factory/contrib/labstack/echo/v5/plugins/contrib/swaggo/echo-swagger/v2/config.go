package swagger

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
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

// ConfigAdd registers the swagger plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable swagger integration")
	config.Add(path+route, "/swagger/*", "define swagger ui route")
}
