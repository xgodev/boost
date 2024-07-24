package http

import (
	cloudevents "github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = cloudevents.Root + ".http"
	port        = root + ".port"
	path        = root + ".path"
	PluginsRoot = root + ".plugins"
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
