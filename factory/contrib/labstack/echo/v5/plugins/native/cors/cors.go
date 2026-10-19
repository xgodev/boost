package cors

import (
	"context"
	"github.com/labstack/echo/v5/middleware"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new cors plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewCors()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Cors represents cors plugin for echo server.
type Cors struct {
	options *Options
}

// NewCorsWithOptions returns a new cors plugin with options.
func NewCorsWithOptions(options *Options) *Cors {
	return &Cors{options: options}
}

// NewCorsWithConfigPath returns a new cors plugin with options from config path.
func NewCorsWithConfigPath(path string) (*Cors, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewCorsWithOptions(o), nil
}

// NewCors returns a new cors plugin with default options.
func NewCors() (*Cors, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewCorsWithOptions(o), nil
}

// Register registers this cors plugin for echo server.
func (i *Cors) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling cors middleware in echo")

	mw, err := middleware.CORSConfig{
		AllowOrigins:     i.options.Allowed.Origins,
		AllowMethods:     i.options.Allowed.Methods,
		AllowHeaders:     i.options.Allowed.Headers,
		AllowCredentials: i.options.Allowed.Credentials,
		ExposeHeaders:    i.options.Exposed.Headers,
		MaxAge:           i.options.MaxAge,
	}.ToMiddleware()
	if err != nil {
		return err
	}
	server.Use(mw)

	logger.Debug("cors middleware successfully enabled in echo")

	return nil
}
