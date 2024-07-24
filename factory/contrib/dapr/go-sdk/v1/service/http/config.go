package http

import (
	"github.com/xgodev/boost/factory/contrib/dapr/go-sdk/v1/service"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = service.Root + ".http"
	address = ".address"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+address, ":8080", ":8080, 0.0.0.0:8080, 10.1.1.1:8080")
}
