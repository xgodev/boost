package health

import (
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v9"
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

// ConfigAdd registers the health plugin keys under path.
func ConfigAdd(path string) {
	config.Add(path+name, "elasticsearch", "health name")
	config.Add(path+description, "default connection", "define health description")
	config.Add(path+required, true, "define if the health is required")
	config.Add(path+enabled, true, "enable/disable health")
}
