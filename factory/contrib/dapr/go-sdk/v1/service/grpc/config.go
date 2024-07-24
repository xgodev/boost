package grpc

import (
	"github.com/xgodev/boost/factory/contrib/dapr/go-sdk/v1/service"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = service.Root + ".grpc"
	address = ".address"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+address, ":50001", ":50001, 0.0.0.0:50001, 10.1.1.1:50001")
}
