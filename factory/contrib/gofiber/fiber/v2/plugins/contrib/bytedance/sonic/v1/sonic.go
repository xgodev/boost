package sonic

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new bytedance/sonic encoder plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewSonic()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Sonic represents bytedance/sonic encoder plugin for fiber server.
type Sonic struct {
	options *Options
}

// NewSonicWithOptions returns a new bytedance/sonic encoder plugin with options.
func NewSonicWithOptions(options *Options) *Sonic {
	return &Sonic{options: options}
}

// NewSonicWithConfigPath returns a new bytedance/sonic encoder plugin with options from config path.
func NewSonicWithConfigPath(path string) (*Sonic, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewSonicWithOptions(o), nil
}

// NewSonic returns a new bytedance/sonic encoder plugin with default options.
func NewSonic() (*Sonic, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewSonicWithOptions(o), nil
}

// Register registers this bytedance/sonic encoder plugin for fiber server.
func (i *Sonic) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling bytedance/sonic encoder in fiber")

	server.Config().JSONEncoder = sonic.Marshal
	server.Config().JSONDecoder = sonic.Unmarshal

	logger.Debug("bytedance/sonic encoder successfully enabled in fiber")

	return nil
}
