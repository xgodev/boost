package prometheus

import (
	"github.com/xgodev/boost/factory/contrib/afex/hystrix-go/v0"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root      = hystrix.PluginsRoot + ".prometheus"
	namespace = ".namespace"
	labels    = ".labels"
	enabled   = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+namespace, "hystrix", "defines hystrix prometheus namespace")
	config.Add(path+labels, map[string]interface{}{}, "defines hystrix prometheus labels")
	config.Add(path+enabled, true, "enabled/disable hystrix prometheus")
}
