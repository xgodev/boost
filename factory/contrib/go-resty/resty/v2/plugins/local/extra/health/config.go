package health

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
)

const (
	root        = resty.PluginsRoot + ".health"
	name        = ".name"
	host        = ".host"
	endpoint    = ".endpoint"
	description = ".description"
	required    = ".required"
	enabled     = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+name, "rest api", "health name")
	config.Add(path+host, "", "health host")
	config.Add(path+endpoint, "/resource-status", "health host")
	config.Add(path+description, "default connection", "define health description")
	config.Add(path+required, true, "define health description")
	config.Add(path+enabled, true, "enable/disable health")
}
