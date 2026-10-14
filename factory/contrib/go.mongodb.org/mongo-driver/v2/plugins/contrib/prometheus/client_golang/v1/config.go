package prometheus

import (
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = resty.PluginsRoot + ".prometheus"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable prometheus integration")
}
