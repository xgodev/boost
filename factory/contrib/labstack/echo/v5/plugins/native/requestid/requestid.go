package requestid

import (
	"context"
	"github.com/labstack/echo/v5/middleware"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new request id plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewRequestID()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// RequestID represents request id plugin for echo server.
type RequestID struct {
	options *Options
}

// NewRequestIDWithOptions returns a new request id plugin with options.
func NewRequestIDWithOptions(options *Options) *RequestID {
	return &RequestID{options: options}
}

// NewRequestIDWithConfigPath returns a new request id plugin with options from config path.
func NewRequestIDWithConfigPath(path string) (*RequestID, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewRequestIDWithOptions(o), nil
}

// NewRequestID returns a new request id plugin with default options.
func NewRequestID() (*RequestID, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewRequestIDWithOptions(o), nil
}

// Register registers this request id plugin for echo server.
func (i *RequestID) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling request id middleware in echo")

	server.Use(middleware.RequestID())

	logger.Debug("request id middleware successfully enabled in echo")

	return nil
}
