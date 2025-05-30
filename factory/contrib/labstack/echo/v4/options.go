package echo

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options echo server options.
type Options struct {
	HideBanner   bool
	DisableHTTP2 bool `config:"disableHTTP2"`
	Port         int
	Type         string
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
	Json struct {
		Pretty struct {
			Enabled bool
		}
	}
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
