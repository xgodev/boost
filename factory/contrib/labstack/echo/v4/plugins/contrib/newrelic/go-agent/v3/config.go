package newrelic

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root                       = echo.PluginsRoot + ".newrelic"
	enabled                    = ".enabled"
	middlewareRoot             = ".middlewares"
	middlewareRequestIDEnabled = middlewareRoot + ".requestId.enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable newrelic integration")
	config.Add(path+middlewareRequestIDEnabled, true, "enable/disable request id middleware")
}
