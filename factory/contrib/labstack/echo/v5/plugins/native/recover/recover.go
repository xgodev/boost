package recover

import (
	"context"
	"github.com/labstack/echo/v5/middleware"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new recover plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewRecover()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Recover represents recover plugin for echo server.
type Recover struct {
	options *Options
}

// NewRecoverWithOptions returns a new recover plugin with options.
func NewRecoverWithOptions(options *Options) *Recover {
	return &Recover{options: options}
}

// NewRecoverWithConfigPath returns a new recover plugin with options from config path.
func NewRecoverWithConfigPath(path string) (*Recover, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewRecoverWithOptions(o), nil
}

// NewRecover returns a new recover plugin with default options.
func NewRecover() (*Recover, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewRecoverWithOptions(o), nil
}

// Register registers this recover plugin for echo server.
func (i *Recover) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling recover middleware in echo")

	server.Use(middleware.Recover())

	logger.Debug("recover middleware successfully enabled in echo")

	return nil
}
