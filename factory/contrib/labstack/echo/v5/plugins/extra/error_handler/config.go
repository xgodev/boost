package error_handler

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = echo.PluginsRoot + ".errorHandler"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the error handler plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable custom error handler")
}
