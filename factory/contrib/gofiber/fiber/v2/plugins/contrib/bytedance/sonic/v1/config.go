package sonic

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = fiber.PluginsRoot + ".bytedanceSonic"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the bytedance/sonic encoder plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable bytedance/sonic encoder")
}
