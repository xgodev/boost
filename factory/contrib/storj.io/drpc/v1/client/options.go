package client

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

// Options drpc client options.
type Options struct {
	Host    string
	Port    int
	Timeout time.Duration
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	ConfigAdd(path)
	config.Load()
	return config.NewOptionsWithPath[Options](root, path)
}
