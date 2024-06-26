package health

import (
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = elasticsearch.PluginsRoot + ".health"
	name        = ".name"
	description = ".description"
	required    = ".required"
	enabled     = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+name, "elasticsearch", "health name")
	config.Add(path+description, "default connection", "define health description")
	config.Add(path+required, true, "define health description")
	config.Add(path+enabled, true, "enable/disable health")
}
