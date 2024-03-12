package prometheus

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
)

const (
	root    = fiber.PluginsRoot + ".prometheus"
	enabled = ".enabled"
	route   = ".route"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable prometheus integration")
	config.Add(path+route, "/metrics", "define prometheus metrics url")
}
