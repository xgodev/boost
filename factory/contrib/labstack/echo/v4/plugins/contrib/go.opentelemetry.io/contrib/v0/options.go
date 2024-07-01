package contrib // import "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/contrib/opentelemetry/otelecho.v1"

import (
	"github.com/xgodev/boost/wrapper/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// Options represents the opentelemetry plugin for echo server options.
type Options struct {
	Enabled        bool
	TracingOptions []otelecho.Option
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
