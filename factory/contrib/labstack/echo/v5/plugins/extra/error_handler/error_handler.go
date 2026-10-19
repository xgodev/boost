package error_handler

import (
	"context"
	e "github.com/labstack/echo/v5"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new error handler plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewErrorHandler()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// ErrorHandler represents error handler plugin for echo server.
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

// Register registers this error handler plugin for echo server.
func (i *ErrorHandler) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling error handler handler in echo")

	server.Instance().HTTPErrorHandler = errorHandler(server)

	logger.Debug("error handler handler successfully enabled in echo")

	return nil
}

// errorHandler answers with JSON for a REST server and plain text otherwise,
// leaving alone a response a handler already started writing.
func errorHandler(server *echo.Server) e.HTTPErrorHandler {
	return func(c *e.Context, err error) {
		if res, uErr := e.UnwrapResponse(c.Response()); uErr == nil && res.Committed {
			return
		}
		if server.Options().Type != "REST" {
			echo.ErrorHandlerString(c, err)
			return
		}
		echo.ErrorHandlerJSON(c, err)
	}
}
