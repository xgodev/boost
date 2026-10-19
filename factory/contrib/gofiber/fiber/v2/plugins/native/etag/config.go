package etag

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = fiber.PluginsRoot + ".etag"
	enabled = ".enabled"
	weak    = ".weak"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the etag middleware plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable etag middleware")
	config.Add(path+weak, false, "generate weak etags")
}
