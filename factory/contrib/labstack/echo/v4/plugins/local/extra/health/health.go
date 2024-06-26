package health

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"

	e "github.com/labstack/echo/v4"
	response "github.com/xgodev/boost/model/restresponse"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new health checker plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	o, err := NewOptions()
	if err != nil {
		return nil
	}
	h := NewHealthWithOptions(o)
	return h.Register(ctx, server)
}

// Health represents health checker plugin for echo server.
type Health struct {
	options *Options
}

// NewHealthWithOptions returns a new health checker plugin with options.
func NewHealthWithOptions(options *Options) *Health {
	return &Health{options: options}
}

// NewHealthWithConfigPath returns a new health checker plugin with options from config path.
func NewHealthWithConfigPath(path string) (*Health, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// NewHealth returns a new health checker plugin with default options.
func NewHealth() *Health {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewHealthWithOptions(o)
}

// Register registers this health checker plugin for echo server.
func (i *Health) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	healthRoute := i.options.Route

	logger.Tracef("configuring health router on %s in echo", healthRoute)

	server.GET(healthRoute, handler)

	logger.Debugf("health router configured on %s in echo", healthRoute)

	return nil
}

func handler(c e.Context) error {

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	resp, httpCode := response.NewHealth(ctx)

	return c.JSON(httpCode, resp)
}
