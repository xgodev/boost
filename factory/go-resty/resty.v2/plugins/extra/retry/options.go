package retry

import (
	"time"

	"github.com/xgodev/boost/factory"
)

// Options represents a retry plugin for resty client options.
type Options struct {
	Enabled     bool
	Count       int
	WaitTime    time.Duration
	MaxWaitTime time.Duration
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return factory.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return factory.NewOptionsWithPath[Options](root, path)
}
