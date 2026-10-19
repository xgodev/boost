package error_handler

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = fiber.PluginsRoot + ".errorHandler"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the error handler plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable error handler")
}
