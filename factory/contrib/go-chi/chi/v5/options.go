package chi

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options chi server options.
type Options struct {
	DisableHTTP2 bool `config:"disableHTTP2"`
	Port         int
	Protocol     string
	TLS          struct {
		Enabled bool
		Type    string
		Auto    struct {
			Host string
		}
		File struct {
			Cert string
			Key  string
		}
	} `config:"tls"`
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	ConfigAdd(path)
	config.Load()
	return config.NewOptionsWithPath[Options](root, path)
}
