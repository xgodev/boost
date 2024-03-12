package datadog

import (
	"github.com/xgodev/boost/config"
)

const (
	root    = echo.PluginsRoot + ".datadog"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable datadog middleware")
}
