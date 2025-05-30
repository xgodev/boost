package sql

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

// Options represents a godror options.
type Options struct {
	ConnMaxLifetime time.Duration
	ConnMaxIdletime time.Duration
	MaxIdletime     int
	MaxOpenConns    int
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
