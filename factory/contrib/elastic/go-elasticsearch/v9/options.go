package elasticsearch

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

// Options are the client settings read from boost.factory.elasticsearch.
type Options struct {
	Addresses             []string
	Username              string
	Password              string
	CloudID               string `config:"cloudID"`
	APIKey                string `config:"APIKey"`
	CACert                string `config:"CACert"`
	RetryOnStatus         []int
	DisableRetry          bool
	MaxRetries            int
	DiscoverNodesOnStart  bool
	DiscoverNodesInterval time.Duration
	EnableMetrics         bool
	EnableDebugLogger     bool
	RetryBackoff          time.Duration
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath returns the root options overridden by the keys under path.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
