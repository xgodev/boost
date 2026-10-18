package client

import (
	"time"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
)

// Transports accepted by Options.Transport (case-insensitive).
const (
	TransportHTTP    = "HTTP"
	TransportCOMMAND = "COMMAND"
)

// Options configures the MCP client and the server it connects to.
type Options struct {
	Name       string
	Version    string
	Transport  string
	Endpoint   string
	MaxRetries int
	Command    string
	Args       []string
	KeepAlive  time.Duration
}

// NewOptions returns options from the default config path.
func NewOptions() (*Options, error) {
	return withDefaults(config.NewOptionsWithPath[Options](root))
}

// NewOptionsWithPath returns options from path, falling back to the root defaults.
func NewOptionsWithPath(path string) (*Options, error) {
	return withDefaults(config.NewOptionsWithPath[Options](root, path))
}

func withDefaults(o *Options, err error) (*Options, error) {
	if err != nil {
		return nil, err
	}
	if o.Name == "" {
		o.Name = boost.ApplicationName()
	}
	return o, nil
}
