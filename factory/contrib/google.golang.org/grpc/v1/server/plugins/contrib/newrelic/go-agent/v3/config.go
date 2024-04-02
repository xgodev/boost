package newrelic

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
)

const (
	root    = server.PluginsRoot + ".newrelic"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable newrelic")
}
