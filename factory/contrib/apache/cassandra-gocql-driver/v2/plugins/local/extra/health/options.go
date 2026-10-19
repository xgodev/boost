package health

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options are the health check settings, read from
// boost.factory.gocql.plugins.health.*.
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

// NewOptionsWithPath returns the options of path, over the root defaults.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
