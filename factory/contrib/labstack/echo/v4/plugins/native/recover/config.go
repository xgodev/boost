package recover

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = echo.PluginsRoot + ".recover"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable recover middleware")
}
