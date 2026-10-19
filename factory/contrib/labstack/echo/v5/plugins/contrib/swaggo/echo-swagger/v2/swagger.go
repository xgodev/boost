// Package swagger serves the Swagger UI and the swag docs on an echo v5
// server, through github.com/swaggo/echo-swagger/v2.
package swagger

import (
	"context"

	eswagger "github.com/swaggo/echo-swagger/v2"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new swagger plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewSwagger()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Swagger represents swagger plugin for echo server.
type Swagger struct {
	options *Options
}

// NewSwaggerWithOptions returns a new swagger plugin with options.
func NewSwaggerWithOptions(options *Options) *Swagger {
	return &Swagger{options: options}
}

// NewSwaggerWithConfigPath returns a new swagger plugin with options from config path.
func NewSwaggerWithConfigPath(path string) (*Swagger, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewSwaggerWithOptions(o), nil
}

// NewSwagger returns a new swagger plugin with default options.
func NewSwagger() (*Swagger, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewSwaggerWithOptions(o), nil
}

// Register registers this swagger plugin for echo server.
func (i *Swagger) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Tracef("configuring swagger router on %s in echo", i.options.Route)

	server.GET(i.options.Route, eswagger.WrapHandler)

	logger.Debugf("swagger router configured on %s in echo", i.options.Route)

	return nil
}
