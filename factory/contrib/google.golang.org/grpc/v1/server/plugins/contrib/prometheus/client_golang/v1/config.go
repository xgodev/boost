package prometheus

import (
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = server.PluginsRoot + ".prometheus"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable prometheus")
}
