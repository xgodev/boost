package otelsql

import (
	"github.com/xgodev/boost/factory/contrib/godror/godror/v0"
	"github.com/xgodev/boost/wrapper/config"
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
