package client

import (
	"github.com/xgodev/boost/factory"
)

// Options kubernetes client set options.
type Options struct {
	Type              string
	KubeConfigPath    string
	KubeConfigContext string
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return factory.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return factory.NewOptionsWithPath[Options](root, path)
}
