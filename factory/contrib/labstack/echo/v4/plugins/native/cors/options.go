package cors

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options CORS plugin for echo server options.
type Options struct {
	Enabled bool
	Allowed struct {
		Origins     []string
		Headers     []string
		Methods     []string
		Credentials bool
	}
	Exposed struct {
		Headers []string
	}
	MaxAge int
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
