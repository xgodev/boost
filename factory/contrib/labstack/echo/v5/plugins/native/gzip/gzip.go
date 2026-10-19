package gzip

import (
	"context"
	"github.com/labstack/echo/v5/middleware"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new gzip plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewGzip()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Gzip represents gzip plugin for echo server.
type Gzip struct {
	options *Options
}

// NewGzipWithOptions returns a new gzip plugin with options.
func NewGzipWithOptions(options *Options) *Gzip {
	return &Gzip{options: options}
}

// NewGzipWithConfigPath returns a new gzip plugin with options from config path.
func NewGzipWithConfigPath(path string) (*Gzip, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewGzipWithOptions(o), nil
}

// NewGzip returns a new gzip plugin with default options.
func NewGzip() (*Gzip, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewGzipWithOptions(o), nil
}

// Register registers this gzip plugin for echo server.
func (i *Gzip) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling gzip middleware in echo")

	mw, err := middleware.GzipConfig{Level: i.options.Level}.ToMiddleware()
	if err != nil {
		return err
	}
	server.Use(mw)

	logger.Debug("gzip middleware successfully enabled in echo")

	return nil
}
