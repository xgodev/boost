package newrelic

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
)

const (
	root    = fiber.PluginsRoot + ".newrelic"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable newrelic middleware")
}
