package compress

import (
	"context"

	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new compress middleware plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewCompress()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Compress represents compress middleware plugin for fiber server.
type Compress struct {
	options *Options
}

// NewCompressWithOptions returns a new compress middleware plugin with options.
func NewCompressWithOptions(options *Options) *Compress {
	return &Compress{options: options}
}

// NewCompressWithConfigPath returns a new compress middleware plugin with options from config path.
func NewCompressWithConfigPath(path string) (*Compress, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewCompressWithOptions(o), nil
}

// NewCompress returns a new compress middleware plugin with default options.
func NewCompress() (*Compress, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewCompressWithOptions(o), nil
}

// Register registers this compress middleware plugin for fiber server.
func (i *Compress) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling compress middleware in fiber")

	server.Use(compress.New(compress.Config{Level: compress.Level(i.options.Level)}))

	logger.Debug("compress middleware successfully enabled in fiber")

	return nil
}
