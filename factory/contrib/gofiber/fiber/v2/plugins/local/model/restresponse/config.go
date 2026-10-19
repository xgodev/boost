package restresponse

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = fiber.PluginsRoot + ".status"
	enabled = ".enabled"
	route   = ".route"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the resource status route plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable resource status route")
	config.Add(path+route, "/resource-status", "define resource status url")
}
