package otel

import "github.com/xgodev/boost/wrapper/config"

// Options configures the otel plugin.
type Options struct {
	Enabled bool
}

// NewOptions returns options from the default config path.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath returns options from path, falling back to the root defaults.
func NewOptionsWithPath(path string) (*Options, error) {
	return config.NewOptionsWithPath[Options](root, path)
}
