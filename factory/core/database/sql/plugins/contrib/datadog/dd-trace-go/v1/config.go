package datadog

import (
	"github.com/xgodev/boost/factory/core/database/sql"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = sql.PluginsRoot + ".datadog"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable datadog integration")
}
