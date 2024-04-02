package compressor

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
)

const (
	root  = server.PluginsRoot + ".compressor"
	level = root + ".level"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+level, -1, "sets gzip level")
}
