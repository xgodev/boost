package log

import (
	ginats "github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = ginats.PluginsRoot + ".log"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable logger")
}
