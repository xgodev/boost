package error_handler

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new error handler plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewErrorHandler()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// ErrorHandler represents error handler plugin for fiber server.
type ErrorHandler struct {
	options *Options
}

// NewErrorHandlerWithOptions returns a new error handler plugin with options.
func NewErrorHandlerWithOptions(options *Options) *ErrorHandler {
	return &ErrorHandler{options: options}
}

// NewErrorHandlerWithConfigPath returns a new error handler plugin with options from config path.
func NewErrorHandlerWithConfigPath(path string) (*ErrorHandler, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewErrorHandlerWithOptions(o), nil
}

// NewErrorHandler returns a new error handler plugin with default options.
func NewErrorHandler() (*ErrorHandler, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewErrorHandlerWithOptions(o), nil
}

// Register registers this error handler plugin for fiber server: a REST
// server answers errors as JSON, any other as plain text.
func (i *ErrorHandler) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling error handler in fiber")

	if server.ServerOptions().Type == "REST" {
		server.Config().ErrorHandler = fiber.ErrorHandlerJSON
	} else {
		server.Config().ErrorHandler = fiber.ErrorHandlerString
	}

	logger.Debug("error handler successfully enabled in fiber")

	return nil
}
