package semaphore

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = echo.PluginsRoot + ".semaphore"
	enabled = ".enabled"
	limit   = ".limit"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable semaphore middleware")
	config.Add(path+limit, int64(10000), "defines numbers for concurrent connections")
}
