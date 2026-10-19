package pprof

import (
	"context"
	"github.com/labstack/echo-contrib/v5/pprof"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new pprof plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewPprof()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Pprof represents pprof plugin for echo server.
type Pprof struct {
	options *Options
}

// NewPprofWithOptions returns a new pprof plugin with options.
func NewPprofWithOptions(options *Options) *Pprof {
	return &Pprof{options: options}
}

// NewPprofWithConfigPath returns a new pprof plugin with options from config path.
func NewPprofWithConfigPath(path string) (*Pprof, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewPprofWithOptions(o), nil
}

// NewPprof returns a new pprof plugin with default options.
func NewPprof() (*Pprof, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewPprofWithOptions(o), nil
}

// Register registers this pprof plugin for echo server.
func (i *Pprof) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling pprof routes in echo")

	pprof.Register(server.Instance(), i.options.Prefix)

	logger.Debug("pprof routes successfully enabled in echo")

	return nil
}
