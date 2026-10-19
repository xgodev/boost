package json

import (
	"context"

	gojson "github.com/goccy/go-json"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new goccy/go-json encoder plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewJSON()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// JSON represents goccy/go-json encoder plugin for fiber server.
type JSON struct {
	options *Options
}

// NewJSONWithOptions returns a new goccy/go-json encoder plugin with options.
func NewJSONWithOptions(options *Options) *JSON {
	return &JSON{options: options}
}

// NewJSONWithConfigPath returns a new goccy/go-json encoder plugin with options from config path.
func NewJSONWithConfigPath(path string) (*JSON, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewJSONWithOptions(o), nil
}

// NewJSON returns a new goccy/go-json encoder plugin with default options.
func NewJSON() (*JSON, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewJSONWithOptions(o), nil
}

// Register registers this goccy/go-json encoder plugin for fiber server.
func (i *JSON) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling goccy/go-json encoder in fiber")

	server.Config().JSONEncoder = gojson.Marshal
	server.Config().JSONDecoder = gojson.Unmarshal

	logger.Debug("goccy/go-json encoder successfully enabled in fiber")

	return nil
}
