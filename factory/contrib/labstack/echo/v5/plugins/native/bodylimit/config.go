package bodylimit

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
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

// ConfigAdd registers the body limit plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable body limit middleware")
	config.Add(path+size, "8M", "body limit size, e.g. 512K, 8M, 1G")
}
