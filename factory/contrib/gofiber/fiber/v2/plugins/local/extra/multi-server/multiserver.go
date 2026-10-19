package multiserver

import (
	"context"

	f "github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new multi server check route plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewMultiServer()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// MultiServer represents multi server check route plugin for fiber server.
type MultiServer struct {
	options *Options
}

// NewMultiServerWithOptions returns a new multi server check route plugin with options.
func NewMultiServerWithOptions(options *Options) *MultiServer {
	return &MultiServer{options: options}
}

// NewMultiServerWithConfigPath returns a new multi server check route plugin with options from config path.
func NewMultiServerWithConfigPath(path string) (*MultiServer, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewMultiServerWithOptions(o), nil
}

// NewMultiServer returns a new multi server check route plugin with default options.
func NewMultiServer() (*MultiServer, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewMultiServerWithOptions(o), nil
}

// Register registers this multi server check route plugin for fiber server.
func (i *MultiServer) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling multi server check route in fiber")

	server.Get(i.options.Route, handler)
	logger.Debugf("multi server check route configured on %s in fiber", i.options.Route)

	logger.Debug("multi server check route successfully enabled in fiber")

	return nil
}

func handler(c *f.Ctx) error {
	if err := multiserver.Check(c.UserContext()); err != nil {
		return c.Status(f.StatusServiceUnavailable).SendString("Service Unavailable")
	}
	return c.SendString("OK")
}
