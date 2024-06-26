package bodylimit

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = echo.PluginsRoot + ".bodyLimit"
	enabled = ".enabled"
	size    = ".size"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(enabled, true, "enable/disable body limit middleware")
	config.Add(size, "8M", "body limit size")
}
