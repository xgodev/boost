package stripslashes

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new strip slashes middleware plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	p, err := NewStripSlashes()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// StripSlashes represents strip slashes middleware plugin for chi server.
type StripSlashes struct {
	options *Options
}

// NewStripSlashesWithOptions returns a new strip slashes middleware plugin with options.
func NewStripSlashesWithOptions(options *Options) *StripSlashes {
	return &StripSlashes{options: options}
}

// NewStripSlashesWithConfigPath returns a new strip slashes middleware plugin with options from config path.
func NewStripSlashesWithConfigPath(path string) (*StripSlashes, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewStripSlashesWithOptions(o), nil
}

// NewStripSlashes returns a new strip slashes middleware plugin with default options.
func NewStripSlashes() (*StripSlashes, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewStripSlashesWithOptions(o), nil
}

// Register registers this strip slashes middleware plugin for chi server.
func (i *StripSlashes) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling strip slashes middleware in chi")

	server.Use(middleware.StripSlashes)

	logger.Debug("strip slashes middleware successfully enabled in chi")

	return nil
}
