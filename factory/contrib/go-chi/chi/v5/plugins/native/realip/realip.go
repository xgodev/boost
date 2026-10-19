package realip

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new real ip middleware plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	p, err := NewRealIP()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// RealIP represents real ip middleware plugin for chi server.
type RealIP struct {
	options *Options
}

// NewRealIPWithOptions returns a new real ip middleware plugin with options.
func NewRealIPWithOptions(options *Options) *RealIP {
	return &RealIP{options: options}
}

// NewRealIPWithConfigPath returns a new real ip middleware plugin with options from config path.
func NewRealIPWithConfigPath(path string) (*RealIP, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewRealIPWithOptions(o), nil
}

// NewRealIP returns a new real ip middleware plugin with default options.
func NewRealIP() (*RealIP, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewRealIPWithOptions(o), nil
}

// Register registers this real ip middleware plugin for chi server.
func (i *RealIP) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling real ip middleware in chi")

	server.Use(middleware.RealIP)

	logger.Debug("real ip middleware successfully enabled in chi")

	return nil
}
