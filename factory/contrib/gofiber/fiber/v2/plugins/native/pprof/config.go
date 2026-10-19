package pprof

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = fiber.PluginsRoot + ".pprof"
	enabled = ".enabled"
	prefix  = ".prefix"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the pprof middleware plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, false, "enable/disable pprof middleware")
	config.Add(path+prefix, "", "prefix of the /debug/pprof routes")
}
