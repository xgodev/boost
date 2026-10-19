package etag

import (
	"context"

	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new etag middleware plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewETag()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// ETag represents etag middleware plugin for fiber server.
type ETag struct {
	options *Options
}

// NewETagWithOptions returns a new etag middleware plugin with options.
func NewETagWithOptions(options *Options) *ETag {
	return &ETag{options: options}
}

// NewETagWithConfigPath returns a new etag middleware plugin with options from config path.
func NewETagWithConfigPath(path string) (*ETag, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewETagWithOptions(o), nil
}

// NewETag returns a new etag middleware plugin with default options.
func NewETag() (*ETag, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewETagWithOptions(o), nil
}

// Register registers this etag middleware plugin for fiber server.
func (i *ETag) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling etag middleware in fiber")

	server.Use(etag.New(etag.Config{Weak: i.options.Weak}))

	logger.Debug("etag middleware successfully enabled in fiber")

	return nil
}
