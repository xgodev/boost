package datadog

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = echo.PluginsRoot + ".datadog"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the datadog plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable datadog middleware")
}
