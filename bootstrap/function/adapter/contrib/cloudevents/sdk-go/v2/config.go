package cloudevents

import (
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root = function.Root + ".cloudevents"
	port = root + ".port"
	path = root + ".path"
)

func init() {
	config.Add(port, 8080, "sets cloudvents port")
	config.Add(path, "/", "sets cloudvents path")
}

func Port() int {
	return config.Int(port)
}

func Path() string {
	return config.String(path)
}
