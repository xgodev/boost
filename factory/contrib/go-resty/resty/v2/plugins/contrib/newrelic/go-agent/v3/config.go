package newrelic

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
)

const (
	root    = resty.PluginsRoot + ".newrelic"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable newrelic integration")
}
