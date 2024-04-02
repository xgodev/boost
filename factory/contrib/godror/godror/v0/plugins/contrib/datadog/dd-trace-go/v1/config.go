package datadog

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/godror/godror/v0"
)

const (
	root    = godror.PluginsRoot + ".datadog"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable datadog integration")
}
