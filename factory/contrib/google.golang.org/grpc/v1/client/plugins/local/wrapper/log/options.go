package log

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options logger plugin for grpc client options.
type Options struct {
	Enabled bool
	Level   string
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
