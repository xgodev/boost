package pprof

import (
	"context"

	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new pprof middleware plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewPprof()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Pprof represents pprof middleware plugin for fiber server.
type Pprof struct {
	options *Options
}

// NewPprofWithOptions returns a new pprof middleware plugin with options.
func NewPprofWithOptions(options *Options) *Pprof {
	return &Pprof{options: options}
}

// NewPprofWithConfigPath returns a new pprof middleware plugin with options from config path.
func NewPprofWithConfigPath(path string) (*Pprof, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewPprofWithOptions(o), nil
}

// NewPprof returns a new pprof middleware plugin with default options.
func NewPprof() (*Pprof, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewPprofWithOptions(o), nil
}

// Register registers this pprof middleware plugin for fiber server.
func (i *Pprof) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling pprof middleware in fiber")

	server.Use(pprof.New(pprof.Config{Prefix: i.options.Prefix}))

	logger.Debug("pprof middleware successfully enabled in fiber")

	return nil
}
