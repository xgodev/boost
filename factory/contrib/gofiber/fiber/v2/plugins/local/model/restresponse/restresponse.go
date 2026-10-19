package restresponse

import (
	"context"

	f "github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	response "github.com/xgodev/boost/model/restresponse"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new resource status route plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewStatus()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Status represents resource status route plugin for fiber server.
type Status struct {
	options *Options
}

// NewStatusWithOptions returns a new resource status route plugin with options.
func NewStatusWithOptions(options *Options) *Status {
	return &Status{options: options}
}

// NewStatusWithConfigPath returns a new resource status route plugin with options from config path.
func NewStatusWithConfigPath(path string) (*Status, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewStatusWithOptions(o), nil
}

// NewStatus returns a new resource status route plugin with default options.
func NewStatus() (*Status, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewStatusWithOptions(o), nil
}

// Register registers this resource status route plugin for fiber server.
func (i *Status) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling resource status route in fiber")

	server.Get(i.options.Route, handler)
	logger.Debugf("resource status route configured on %s in fiber", i.options.Route)

	logger.Debug("resource status route successfully enabled in fiber")

	return nil
}

func handler(c *f.Ctx) error {
	return c.JSON(response.NewResourceStatus())
}
