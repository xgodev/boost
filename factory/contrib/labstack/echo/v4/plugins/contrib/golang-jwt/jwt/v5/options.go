package jwt

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

// Options jwt plugin for echo server options.
type Options struct {
	Enabled   bool
	JwksUrl   string
	Issuer    string
	Audience  string
	SkipPaths []string
	Leeway    time.Duration
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
