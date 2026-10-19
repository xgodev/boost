package prometheus

import (
	"context"
	"strconv"
	"sync"
	"time"

	f "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new prometheus integration plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewPrometheus()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Prometheus represents prometheus integration plugin for fiber server.
type Prometheus struct {
	options *Options
}

// NewPrometheusWithOptions returns a new prometheus integration plugin with options.
func NewPrometheusWithOptions(options *Options) *Prometheus {
	return &Prometheus{options: options}
}

// NewPrometheusWithConfigPath returns a new prometheus integration plugin with options from config path.
func NewPrometheusWithConfigPath(path string) (*Prometheus, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewPrometheusWithOptions(o), nil
}

// NewPrometheus returns a new prometheus integration plugin with default options.
func NewPrometheus() (*Prometheus, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewPrometheusWithOptions(o), nil
}

// Register registers this prometheus integration plugin for fiber server.
func (i *Prometheus) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling prometheus integration in fiber")

	if i.options.Collector.Enabled {
		mw, err := metricsMiddleware()
		if err != nil {
			return err
		}
		server.Use(mw)
	}

	if i.options.Route.Enabled {
		server.Get(i.options.Route.Path, adaptor.HTTPHandler(promhttp.Handler()))
		logger.Debugf("prometheus metrics route configured on %s in fiber", i.options.Route.Path)
	}

	logger.Debug("prometheus integration successfully enabled in fiber")

	return nil
}

var (
	collectorOnce sync.Once
	collector     f.Handler
	collectorErr  error
)

// metricsMiddleware registers the fiber collectors in the default prometheus
// registry once per process and shares them: every fiber server in the app
// reports into the same fiber_* metrics, as a second registration would fail.
func metricsMiddleware() (f.Handler, error) {
	collectorOnce.Do(func() {
		collector, collectorErr = newMetricsMiddleware(prometheus.DefaultRegisterer)
	})
	return collector, collectorErr
}

func newMetricsMiddleware(reg prometheus.Registerer) (f.Handler, error) {
	labels := []string{"code", "method", "url"}
	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fiber_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code, method and route.",
	}, labels)
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "fiber_request_duration_seconds",
		Help:    "The HTTP request latencies in seconds, partitioned by status code, method and route.",
		Buckets: prometheus.DefBuckets,
	}, labels)
	for _, col := range []prometheus.Collector{requests, duration} {
		if err := reg.Register(col); err != nil {
			return nil, err
		}
	}

	return func(c *f.Ctx) error {
		start := time.Now()

		// An error from the handlers goes through the app error handler
		// first, so the counted status is the one the client gets.
		if err := c.Next(); err != nil {
			if herr := c.App().ErrorHandler(c, err); herr != nil {
				_ = c.SendStatus(f.StatusInternalServerError) // intentional: the caller is gone when the write fails
			}
		}

		// The route path, not the request path, keeps one series per route
		// instead of one per id.
		values := []string{strconv.Itoa(c.Response().StatusCode()), c.Method(), c.Route().Path}
		requests.WithLabelValues(values...).Inc()
		duration.WithLabelValues(values...).Observe(time.Since(start).Seconds())
		return nil
	}, nil
}
