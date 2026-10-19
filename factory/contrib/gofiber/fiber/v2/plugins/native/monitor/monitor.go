package monitor

import (
	"context"

	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new monitor route plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewMonitor()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Monitor represents monitor route plugin for fiber server.
type Monitor struct {
	options *Options
}

// NewMonitorWithOptions returns a new monitor route plugin with options.
func NewMonitorWithOptions(options *Options) *Monitor {
	return &Monitor{options: options}
}

// NewMonitorWithConfigPath returns a new monitor route plugin with options from config path.
func NewMonitorWithConfigPath(path string) (*Monitor, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewMonitorWithOptions(o), nil
}

// NewMonitor returns a new monitor route plugin with default options.
func NewMonitor() (*Monitor, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewMonitorWithOptions(o), nil
}

// Register registers this monitor route plugin for fiber server.
func (i *Monitor) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling monitor route in fiber")

	server.Get(i.options.Route, monitor.New(monitor.Config{Title: i.options.Title}))
	logger.Debugf("monitor route configured on %s in fiber", i.options.Route)

	logger.Debug("monitor route successfully enabled in fiber")

	return nil
}
