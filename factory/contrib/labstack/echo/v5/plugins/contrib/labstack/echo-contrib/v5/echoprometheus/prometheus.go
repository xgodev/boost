package echoprometheus

import (
	"context"
	"sync"

	"github.com/labstack/echo-contrib/v5/echoprometheus"
	e "github.com/labstack/echo/v5"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new prometheus plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewPrometheus()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Prometheus represents prometheus plugin for echo server.
type Prometheus struct {
	options *Options
}

// NewPrometheusWithOptions returns a new prometheus plugin with options.
func NewPrometheusWithOptions(options *Options) *Prometheus {
	return &Prometheus{options: options}
}

// NewPrometheusWithConfigPath returns a new prometheus plugin with options from config path.
func NewPrometheusWithConfigPath(path string) (*Prometheus, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewPrometheusWithOptions(o), nil
}

// NewPrometheus returns a new prometheus plugin with default options.
func NewPrometheus() (*Prometheus, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewPrometheusWithOptions(o), nil
}

// Register registers this prometheus plugin for echo server.
func (i *Prometheus) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling prometheus integration in echo")

	if i.options.Collector.Enabled {
		mw, err := metricsMiddleware()
		if err != nil {
			return err
		}
		server.Use(mw)
	}

	if i.options.Route.Enabled {
		server.GET(i.options.Route.Path, echoprometheus.NewHandler())
		logger.Debugf("prometheus metrics route configured on %s in echo", i.options.Route.Path)
	}

	logger.Debug("prometheus integration successfully enabled in echo")

	return nil
}

var (
	collectorOnce sync.Once
	collector     e.MiddlewareFunc
	collectorErr  error
)

// metricsMiddleware registers the echo collectors in the default prometheus
// registry once per process and shares them: every echo server in the app
// reports into the same echo_* metrics, as a second registration would fail.
func metricsMiddleware() (e.MiddlewareFunc, error) {
	collectorOnce.Do(func() {
		collector, collectorErr = echoprometheus.MiddlewareConfig{}.ToMiddleware()
	})
	return collector, collectorErr
}
