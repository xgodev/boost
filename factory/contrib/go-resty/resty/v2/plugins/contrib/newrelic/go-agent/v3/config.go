package newrelic

import (
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
	"github.com/xgodev/boost/wrapper/config"
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
