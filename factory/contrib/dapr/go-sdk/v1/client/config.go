package client

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = "boost.factory.dapr"
	address = ".address"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+address, "127.0.0.1:50001", "The cache size will be set to 512KB at minimum")
}
