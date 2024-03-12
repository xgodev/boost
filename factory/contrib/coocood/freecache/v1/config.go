package freecache

import (
	"github.com/xgodev/boost/config"
)

const (
	root      = "boost.factory.freecache"
	cacheSize = ".cacheSize"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+cacheSize, 100*1024*1024, "The cache size will be set to 512KB at minimum")
}
