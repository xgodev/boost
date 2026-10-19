package health

import (
	"context"

	f "github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	response "github.com/xgodev/boost/model/restresponse"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new health route plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewHealth()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Health represents health route plugin for fiber server.
type Health struct {
	options *Options
}

// NewHealthWithOptions returns a new health route plugin with options.
func NewHealthWithOptions(options *Options) *Health {
	return &Health{options: options}
}

// NewHealthWithConfigPath returns a new health route plugin with options from config path.
func NewHealthWithConfigPath(path string) (*Health, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// NewHealth returns a new health route plugin with default options.
func NewHealth() (*Health, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// Register registers this health route plugin for fiber server.
func (i *Health) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling health route in fiber")

	server.Get(i.options.Route, handler)
	logger.Debugf("health route configured on %s in fiber", i.options.Route)

	logger.Debug("health route successfully enabled in fiber")

	return nil
}

func handler(c *f.Ctx) error {
	resp, httpCode := response.NewHealth(c.UserContext())
	return c.Status(httpCode).JSON(resp)
}
