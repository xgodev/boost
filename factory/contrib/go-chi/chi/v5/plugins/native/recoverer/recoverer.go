package recoverer

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new recoverer middleware plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	p, err := NewRecoverer()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Recoverer represents recoverer middleware plugin for chi server.
type Recoverer struct {
	options *Options
}

// NewRecovererWithOptions returns a new recoverer middleware plugin with options.
func NewRecovererWithOptions(options *Options) *Recoverer {
	return &Recoverer{options: options}
}

// NewRecovererWithConfigPath returns a new recoverer middleware plugin with options from config path.
func NewRecovererWithConfigPath(path string) (*Recoverer, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewRecovererWithOptions(o), nil
}

// NewRecoverer returns a new recoverer middleware plugin with default options.
func NewRecoverer() (*Recoverer, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewRecovererWithOptions(o), nil
}

// Register registers this recoverer middleware plugin for chi server.
func (i *Recoverer) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling recoverer middleware in chi")

	server.Use(middleware.Recoverer)

	logger.Debug("recoverer middleware successfully enabled in chi")

	return nil
}
