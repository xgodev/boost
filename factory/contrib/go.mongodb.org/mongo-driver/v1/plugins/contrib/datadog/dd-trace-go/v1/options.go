package datadog

import (
	"github.com/xgodev/boost"
)

// Options represents datadog plugin for mongo options.
type Options struct {
	Enabled       bool
	ServiceName   string
	Analytics     bool
	AnalyticsRate float64
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return boost.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return boost.NewOptionsWithPath[Options](root, path)
}
