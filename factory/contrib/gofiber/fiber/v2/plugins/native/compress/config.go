package compress

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = fiber.PluginsRoot + ".compress"
	enabled = ".enabled"
	level   = ".level"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the compress middleware plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable compress middleware")
	config.Add(path+level, 0, "compression level: -1 disabled, 0 default, 1 best speed, 2 best compression")
}
