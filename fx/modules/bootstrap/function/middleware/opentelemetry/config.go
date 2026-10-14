package opentelemetry

import (
	"github.com/xgodev/boost/bootstrap/function/middleware"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = middleware.Root + ".opentelemetry"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the opentelemetry middleware")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
