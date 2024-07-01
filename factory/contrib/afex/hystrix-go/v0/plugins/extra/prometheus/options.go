package prometheus

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options represents health options.
type Options struct {
	Namespace string
	Labels    map[string]string
	Enabled   bool
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
