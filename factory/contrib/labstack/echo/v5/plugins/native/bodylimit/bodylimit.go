package bodylimit

import (
	"context"
	"github.com/labstack/echo/v5/middleware"
	"github.com/labstack/gommon/bytes"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new body limit plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewBodyLimit()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// BodyLimit represents body limit plugin for echo server.
type BodyLimit struct {
	options *Options
}

// NewBodyLimitWithOptions returns a new body limit plugin with options.
func NewBodyLimitWithOptions(options *Options) *BodyLimit {
	return &BodyLimit{options: options}
}

// NewBodyLimitWithConfigPath returns a new body limit plugin with options from config path.
func NewBodyLimitWithConfigPath(path string) (*BodyLimit, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewBodyLimitWithOptions(o), nil
}

// NewBodyLimit returns a new body limit plugin with default options.
func NewBodyLimit() (*BodyLimit, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewBodyLimitWithOptions(o), nil
}

// Register registers this body limit plugin for echo server.
func (i *BodyLimit) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling body limit middleware in echo")

	limit, err := bytes.Parse(i.options.Size)
	if err != nil {
		return err
	}
	mw, err := middleware.BodyLimitConfig{LimitBytes: limit}.ToMiddleware()
	if err != nil {
		return err
	}
	server.Use(mw)

	logger.Debug("body limit middleware successfully enabled in echo")

	return nil
}
