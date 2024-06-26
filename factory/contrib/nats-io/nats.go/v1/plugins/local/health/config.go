package health

import (
	ginats "github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = ginats.PluginsRoot + ".health"
	name        = ".name"
	description = ".description"
	required    = ".required"
	enabled     = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+name, "nats", "health name")
	config.Add(path+description, "default connection", "define health description")
	config.Add(path+required, true, "define health description")
	config.Add(path+enabled, true, "enable/disable health")
}
