package health

import (
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = client.PluginsRoot + ".health"
	name        = ".name"
	description = ".description"
	required    = ".required"
	enabled     = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+name, "grpc", "health name")
	config.Add(path+description, "default connection", "define health description")
	config.Add(path+required, true, "define health description")
	config.Add(path+enabled, true, "enable/disable health")
}
