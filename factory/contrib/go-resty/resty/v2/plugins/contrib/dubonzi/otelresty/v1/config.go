package otelresty

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
)

const (
	root       = resty.PluginsRoot + ".otel"
	enabled    = ".enabled"
	tracerName = ".tracerName"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable the opentelemetry integration")
	config.Add(path+tracerName, "resty.request", "defines the name of the tracer used to create spans")
}
