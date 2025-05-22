package otel

import (
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = redis.PluginsRoot + ".otel"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable otel integration")
}
