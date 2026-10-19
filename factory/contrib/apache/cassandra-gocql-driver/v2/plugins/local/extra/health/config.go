package health

import (
	gocql "github.com/xgodev/boost/factory/contrib/apache/cassandra-gocql-driver/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = gocql.PluginsRoot + ".health"
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
	config.Add(path+name, "cassandra", "health name")
	config.Add(path+description, "default connection", "define health description")
	config.Add(path+required, true, "define if the health check is required")
	config.Add(path+enabled, true, "enable/disable health")
}
