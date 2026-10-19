package cors

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new cors middleware plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewCors()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Cors represents cors middleware plugin for fiber server.
type Cors struct {
	options *Options
}

// NewCorsWithOptions returns a new cors middleware plugin with options.
func NewCorsWithOptions(options *Options) *Cors {
	return &Cors{options: options}
}

// NewCorsWithConfigPath returns a new cors middleware plugin with options from config path.
func NewCorsWithConfigPath(path string) (*Cors, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewCorsWithOptions(o), nil
}

// NewCors returns a new cors middleware plugin with default options.
func NewCors() (*Cors, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewCorsWithOptions(o), nil
}

// Register registers this cors middleware plugin for fiber server.
func (i *Cors) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling cors middleware in fiber")

	server.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(i.options.Allowed.Origins, ","),
		AllowMethods:     strings.Join(i.options.Allowed.Methods, ","),
		AllowHeaders:     strings.Join(i.options.Allowed.Headers, ","),
		AllowCredentials: i.options.Allowed.Credentials,
		ExposeHeaders:    strings.Join(i.options.Exposed.Headers, ","),
		MaxAge:           i.options.MaxAge,
	}))

	logger.Debug("cors middleware successfully enabled in fiber")

	return nil
}
