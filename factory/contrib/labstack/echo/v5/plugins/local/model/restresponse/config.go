package status

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = echo.PluginsRoot + ".status"
	enabled = ".enabled"
	route   = ".route"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the resource status plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable resource status route")
	config.Add(path+route, "/resource-status", "define resource status url")
}
