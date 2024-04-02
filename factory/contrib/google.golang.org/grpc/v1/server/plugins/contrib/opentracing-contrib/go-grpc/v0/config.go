package opentracing

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
)

const (
	root    = server.PluginsRoot + ".opentracing"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable opentracing")
}
