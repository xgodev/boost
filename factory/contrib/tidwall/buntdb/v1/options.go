package buntdb

import (
	"github.com/xgodev/boost/wrapper/config"
)

type Options struct {
	Path       string
	SyncPolicy int
	AutoShrink struct {
		Percentage int
		MinSize    int
		Disabled   bool
	}
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
