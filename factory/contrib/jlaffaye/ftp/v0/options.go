package ftp

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options ftp server connection options.
type Options struct {
	Addr     string
	User     string
	Password string
	Timeout  int
	Retry    int
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
