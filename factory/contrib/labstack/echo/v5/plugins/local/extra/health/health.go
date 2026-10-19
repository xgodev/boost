package health

import (
	"context"
	e "github.com/labstack/echo/v5"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	response "github.com/xgodev/boost/model/restresponse"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new health plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewHealth()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Health represents health plugin for echo server.
type Health struct {
	options *Options
}

// NewHealthWithOptions returns a new health plugin with options.
func NewHealthWithOptions(options *Options) *Health {
	return &Health{options: options}
}

// NewHealthWithConfigPath returns a new health plugin with options from config path.
func NewHealthWithConfigPath(path string) (*Health, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// NewHealth returns a new health plugin with default options.
func NewHealth() (*Health, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// Register registers this health plugin for echo server.
func (i *Health) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling health route in echo")

	route := i.options.Route
	server.GET(route, handler)
	logger.Debugf("health route configured on %s in echo", route)

	logger.Debug("health route successfully enabled in echo")

	return nil
}

func handler(c *e.Context) error {
	resp, httpCode := response.NewHealth(c.Request().Context())

	return c.JSON(httpCode, resp)
}
