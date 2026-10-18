package server

import (
	"time"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/internal/buildinfo"
	"github.com/xgodev/boost/wrapper/config"
)

// Transports accepted by Options.Transport (case-insensitive).
const (
	TransportHTTP  = "HTTP"
	TransportSTDIO = "STDIO"
)

// Options configures the MCP server.
type Options struct {
	Name           string
	Version        string
	Instructions   string
	Transport      string
	Port           int
	Path           string
	Stateless      bool
	JSONResponse   bool
	SessionTimeout time.Duration
	KeepAlive      time.Duration
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
	if o.Version == "" {
		o.Version = buildinfo.Version()
	}
	return o, nil
}
