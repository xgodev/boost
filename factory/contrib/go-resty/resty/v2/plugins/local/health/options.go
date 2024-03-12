package health

import (
	"github.com/xgodev/boost"
)

// Options represents resty health options.
type Options struct {
	Name        string
	Host        string
	Endpoint    string
	Enabled     bool
	Description string
	Required    bool
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return boost.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return boost.NewOptionsWithPath[Options](root, path)
}
