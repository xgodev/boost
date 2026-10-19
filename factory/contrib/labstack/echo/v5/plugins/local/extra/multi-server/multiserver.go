package multiserver

import (
	"context"
	"net/http"

	e "github.com/labstack/echo/v5"
	"github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new multi server check plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewMultiServer()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// MultiServer represents multi server check plugin for echo server.
type MultiServer struct {
	options *Options
}

// NewMultiServerWithOptions returns a new multi server check plugin with options.
func NewMultiServerWithOptions(options *Options) *MultiServer {
	return &MultiServer{options: options}
}

// NewMultiServerWithConfigPath returns a new multi server check plugin with options from config path.
func NewMultiServerWithConfigPath(path string) (*MultiServer, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewMultiServerWithOptions(o), nil
}

// NewMultiServer returns a new multi server check plugin with default options.
func NewMultiServer() (*MultiServer, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewMultiServerWithOptions(o), nil
}

// Register registers this multi server check plugin for echo server.
func (i *MultiServer) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling multi server check route in echo")

	route := i.options.Route
	server.GET(route, handler)
	logger.Debugf("multi server check route configured on %s in echo", route)

	logger.Debug("multi server check route successfully enabled in echo")

	return nil
}

func handler(c *e.Context) error {
	if err := multiserver.Check(c.Request().Context()); err != nil {
		return c.String(http.StatusServiceUnavailable, "Service Unavailable")
	}
	return c.String(http.StatusOK, "OK")
}
