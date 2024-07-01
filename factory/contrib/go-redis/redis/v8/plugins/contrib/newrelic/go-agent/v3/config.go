package newrelic

import (
	"github.com/xgodev/boost/factory/contrib/go-redis/redis/v8"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = redis.PluginsRoot + ".newrelic"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable newrelic integration")
}
