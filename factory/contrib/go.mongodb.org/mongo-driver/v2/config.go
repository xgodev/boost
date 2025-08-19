package mongo

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = "boost.factory.mongo"
	uri         = ".uri"
	authRoot    = ".auth"
	username    = authRoot + ".username"
	password    = authRoot + ".password"
	PluginsRoot = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+uri, "mongodb://localhost:27017/temp", "define mongodb uri")
	config.Add(path+username, "", "define mongodb username", config.WithHide())
	config.Add(path+password, "", "define mongodb password", config.WithHide())
}
