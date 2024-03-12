package azure

import (
	"github.com/xgodev/boost/config"
)

const (
	root = "faas.azure"
	port = root + ".port"
	name = root + ".name"
)

func init() {
	config.Add(port, "7071", "define http port")
	config.Add(name, "handler", "define name")
}
