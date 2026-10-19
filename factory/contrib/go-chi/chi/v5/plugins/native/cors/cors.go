package cors

import (
	"context"
	"github.com/go-chi/cors"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new cors middleware plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	p, err := NewCors()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Cors represents cors middleware plugin for chi server.
type Cors struct {
	options *Options
}

// NewCorsWithOptions returns a new cors middleware plugin with options.
func NewCorsWithOptions(options *Options) *Cors {
	return &Cors{options: options}
}

// NewCorsWithConfigPath returns a new cors middleware plugin with options from config path.
func NewCorsWithConfigPath(path string) (*Cors, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewCorsWithOptions(o), nil
}

// NewCors returns a new cors middleware plugin with default options.
func NewCors() (*Cors, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewCorsWithOptions(o), nil
}

// Register registers this cors middleware plugin for chi server.
func (i *Cors) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling cors middleware in chi")

	server.Use(cors.Handler(cors.Options{
		AllowedOrigins:   i.options.Allowed.Origins,
		AllowedMethods:   i.options.Allowed.Methods,
		AllowedHeaders:   i.options.Allowed.Headers,
		AllowCredentials: i.options.Allowed.Credentials,
		ExposedHeaders:   i.options.Exposed.Headers,
		MaxAge:           i.options.MaxAge,
	}))

	logger.Debug("cors middleware successfully enabled in chi")

	return nil
}
