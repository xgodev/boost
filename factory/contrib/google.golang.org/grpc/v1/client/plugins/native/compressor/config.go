package compressor

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
)

const (
	root  = client.PluginsRoot + ".compressor"
	level = ".level"
)

func ConfigAdd(path string) {
	config.Add(path+level, -1, "sets gzip level")
}
