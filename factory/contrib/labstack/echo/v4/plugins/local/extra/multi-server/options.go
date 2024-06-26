package multiserver

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options multiserver plugin for echo server options.
type Options struct {
	Enabled bool
	Route   string
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
