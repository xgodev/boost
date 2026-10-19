package status

import (
	"context"
	"net/http"

	e "github.com/labstack/echo/v5"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	response "github.com/xgodev/boost/model/restresponse"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new resource status plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewStatus()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Status represents resource status plugin for echo server.
type Status struct {
	options *Options
}

// NewStatusWithOptions returns a new resource status plugin with options.
func NewStatusWithOptions(options *Options) *Status {
	return &Status{options: options}
}

// NewStatusWithConfigPath returns a new resource status plugin with options from config path.
func NewStatusWithConfigPath(path string) (*Status, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewStatusWithOptions(o), nil
}

// NewStatus returns a new resource status plugin with default options.
func NewStatus() (*Status, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewStatusWithOptions(o), nil
}

// Register registers this resource status plugin for echo server.
func (i *Status) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling resource status route in echo")

	route := i.options.Route
	server.GET(route, handler)
	logger.Debugf("resource status route configured on %s in echo", route)

	logger.Debug("resource status route successfully enabled in echo")

	return nil
}

func handler(c *e.Context) error {
	return c.JSON(http.StatusOK, response.NewResourceStatus())
}
