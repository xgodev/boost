package compressor

import (
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	"github.com/xgodev/boost/wrapper/config"
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
