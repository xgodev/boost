package otel

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root       = echo.PluginsRoot + ".otel"
	enabled    = ".enabled"
	serverName = ".serverName"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the opentelemetry plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable the opentelemetry integration")
	config.Add(path+serverName, "", "server name recorded on spans and metrics; empty uses the request host")
}
