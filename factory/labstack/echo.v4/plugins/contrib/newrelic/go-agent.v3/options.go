package newrelic

import "github.com/xgodev/boost/factory"

// Options newrelic plugin for echo server options.
type Options struct {
	Enabled     bool
	Middlewares struct {
		RequestID struct {
			Enabled bool
		} `config:"requestId"`
	}
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return factory.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return factory.NewOptionsWithPath[Options](root, path)
}
