package pgx

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options represents a godror options.
type Options struct {
	ConnectString string
	Username      string
	Password      string
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
