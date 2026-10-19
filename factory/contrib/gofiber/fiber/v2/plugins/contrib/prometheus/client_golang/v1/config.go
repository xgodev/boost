package prometheus

import (
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root             = fiber.PluginsRoot + ".prometheus"
	enabled          = ".enabled"
	route            = ".route.path"
	routeEnabled     = ".route.enabled"
	collectorEnabled = ".collector.enabled"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the prometheus integration plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable prometheus integration")
	config.Add(path+route, "/metrics", "define prometheus metrics url")
	config.Add(path+routeEnabled, true, "enable/disable prometheus metrics url")
	config.Add(path+collectorEnabled, true, "enable/disable prometheus collector")
}
