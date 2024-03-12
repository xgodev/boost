package opentracing

import (
	"github.com/xgodev/boost/config"
)

const (
	root    = resty.PluginsRoot + ".opentracing"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable opentracing integration")
}
