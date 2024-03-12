package health

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v7"
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
