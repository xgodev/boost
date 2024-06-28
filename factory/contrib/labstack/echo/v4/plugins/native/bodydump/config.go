package bodydump

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = echo.PluginsRoot + ".bodydump"
	enabled = ".enabled"
)

func init() {

}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable body dump middleware")
}
