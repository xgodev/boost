package monitor

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = fiber.PluginsRoot + ".monitor"
	enabled = ".enabled"
	route   = ".route"
	title   = ".title"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the monitor route plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, false, "enable/disable monitor route")
	config.Add(path+route, "/monitor", "define monitor url")
	config.Add(path+title, "fiber monitor", "monitor page title")
}
