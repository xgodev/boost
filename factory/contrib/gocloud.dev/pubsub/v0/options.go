package pubsub

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options represents pubsub client options.
type Options struct {
	Resource string
	Type     string
	Region   string
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
