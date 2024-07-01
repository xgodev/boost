package health

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options struct which represents a health plugin for redis options.
type Options struct {
	Name        string
	Enabled     bool
	Description string
	Required    bool
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
