package requestid

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = fiber.PluginsRoot + ".requestId"
	enabled = ".enabled"
	header  = ".header"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the request id middleware plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable request id middleware")
	config.Add(path+header, "X-Request-ID", "header the request id is read from and written to")
}
