package gzip

import (
	"context"

	"github.com/labstack/echo/v4/middleware"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new gzip plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	o, err := NewOptions()
	if err != nil {
		return nil
	}
	h := NewGzipWithOptions(o)
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
func NewGzip() *Gzip {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewGzipWithOptions(o)
}

// Register registers this gzip plugin for echo server.
func (i *Gzip) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling gzip middleware in echo")

	server.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: middleware.DefaultGzipConfig.Skipper,
		Level:   i.options.Level,
	}))

	logger.Debug("gzip middleware successfully enabled in echo")

	return nil
}
