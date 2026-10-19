package recover

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root       = fiber.PluginsRoot + ".recover"
	enabled    = ".enabled"
	stackTrace = ".stackTrace"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the recover middleware plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable recover middleware")
	config.Add(path+stackTrace, false, "print the stack trace of the panic")
}
