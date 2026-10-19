package requestid

import (
	"context"

	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new request id middleware plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewRequestID()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// RequestID represents request id middleware plugin for fiber server.
type RequestID struct {
	options *Options
}

// NewRequestIDWithOptions returns a new request id middleware plugin with options.
func NewRequestIDWithOptions(options *Options) *RequestID {
	return &RequestID{options: options}
}

// NewRequestIDWithConfigPath returns a new request id middleware plugin with options from config path.
func NewRequestIDWithConfigPath(path string) (*RequestID, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewRequestIDWithOptions(o), nil
}

// NewRequestID returns a new request id middleware plugin with default options.
func NewRequestID() (*RequestID, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewRequestIDWithOptions(o), nil
}

// Register registers this request id middleware plugin for fiber server.
func (i *RequestID) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling request id middleware in fiber")

	server.Use(requestid.New(requestid.Config{Header: i.options.Header}))

	logger.Debug("request id middleware successfully enabled in fiber")

	return nil
}
