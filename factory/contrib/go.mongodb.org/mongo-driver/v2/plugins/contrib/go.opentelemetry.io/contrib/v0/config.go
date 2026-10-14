package contrib // import "github.com/xgodev/boost/factory/go.mongodb.org/mongo-driver.v1/plugins/contrib/opentelemetry/otelmongo.v1

import (
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = mongo.PluginsRoot + ".otel"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable the opentelemetry integration")
}
