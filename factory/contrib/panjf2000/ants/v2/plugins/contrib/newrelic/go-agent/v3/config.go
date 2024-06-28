package newrelic

import (
	"github.com/xgodev/boost/factory/contrib/panjf2000/ants/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = ants.PluginsRoot + ".newrelic"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable newrelic integration")
}
