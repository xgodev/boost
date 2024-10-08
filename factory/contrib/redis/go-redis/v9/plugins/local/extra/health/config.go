package health

import (
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = redis.PluginsRoot + ".health"
	name        = ".name"
	description = ".description"
	required    = ".required"
	enabled     = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+name, "redis", "health name")
	config.Add(path+description, "default connection", "define health description")
	config.Add(path+required, true, "define health description")
	config.Add(path+enabled, true, "enable/disable health")
}
