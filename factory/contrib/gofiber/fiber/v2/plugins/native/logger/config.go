package logger

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = fiber.PluginsRoot + ".logger"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable logger middleware")
}
