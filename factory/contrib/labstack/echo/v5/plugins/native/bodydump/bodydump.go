package bodydump

import (
	"context"
	e "github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new body dump plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewBodyDump()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// BodyDump represents body dump plugin for echo server.
type BodyDump struct {
	options *Options
}

// NewBodyDumpWithOptions returns a new body dump plugin with options.
func NewBodyDumpWithOptions(options *Options) *BodyDump {
	return &BodyDump{options: options}
}

// NewBodyDumpWithConfigPath returns a new body dump plugin with options from config path.
func NewBodyDumpWithConfigPath(path string) (*BodyDump, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewBodyDumpWithOptions(o), nil
}

// NewBodyDump returns a new body dump plugin with default options.
func NewBodyDump() (*BodyDump, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewBodyDumpWithOptions(o), nil
}

// Register registers this body dump plugin for echo server.
func (i *BodyDump) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling body dump middleware in echo")

	server.Use(middleware.BodyDump(bodyDump))

	logger.Debug("body dump middleware successfully enabled in echo")

	return nil
}

func bodyDump(c *e.Context, reqBody []byte, resBody []byte, _ error) {
	logger := log.FromContext(c.Request().Context())
	logger.Info("request body --->")
	logger.Info(string(reqBody))
	logger.Info("response body -->")
	logger.Info(string(resBody))
}
