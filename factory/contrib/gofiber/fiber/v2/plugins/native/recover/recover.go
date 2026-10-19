package recover

import (
	"context"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new recover middleware plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewRecover()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Recover represents recover middleware plugin for fiber server.
type Recover struct {
	options *Options
}

// NewRecoverWithOptions returns a new recover middleware plugin with options.
func NewRecoverWithOptions(options *Options) *Recover {
	return &Recover{options: options}
}

// NewRecoverWithConfigPath returns a new recover middleware plugin with options from config path.
func NewRecoverWithConfigPath(path string) (*Recover, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewRecoverWithOptions(o), nil
}

// NewRecover returns a new recover middleware plugin with default options.
func NewRecover() (*Recover, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewRecoverWithOptions(o), nil
}

// Register registers this recover middleware plugin for fiber server.
func (i *Recover) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling recover middleware in fiber")

	server.Use(recover.New(recover.Config{EnableStackTrace: i.options.StackTrace}))

	logger.Debug("recover middleware successfully enabled in fiber")

	return nil
}
