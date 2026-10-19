package prometheus

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	c "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new prometheus plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	h, err := NewPrometheus()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Prometheus represents prometheus plugin for chi server.
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

// Register registers this prometheus plugin for chi server.
func (i *Prometheus) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling prometheus integration in chi")

	if i.options.Collector.Enabled {
		mw, err := metricsMiddleware()
		if err != nil {
			return err
		}
		server.Use(mw)
	}

	if i.options.Route.Enabled {
		server.Handle(i.options.Route.Path, promhttp.Handler())
		logger.Debugf("prometheus metrics route configured on %s in chi", i.options.Route.Path)
	}

	logger.Debug("prometheus integration successfully enabled in chi")

	return nil
}

var (
	collectorOnce sync.Once
	collector     func(http.Handler) http.Handler
	collectorErr  error
)

// metricsMiddleware registers the chi collectors in the default prometheus
// registry once per process and shares them: every chi server in the app
// reports into the same chi_* metrics, as a second registration would fail.
func metricsMiddleware() (func(http.Handler) http.Handler, error) {
	collectorOnce.Do(func() {
		collector, collectorErr = newMetricsMiddleware(prometheus.DefaultRegisterer)
	})
	return collector, collectorErr
}

func newMetricsMiddleware(reg prometheus.Registerer) (func(http.Handler) http.Handler, error) {
	labels := []string{"code", "method", "url"}
	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "chi_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code, method and route.",
	}, labels)
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "chi_request_duration_seconds",
		Help:    "The HTTP request latencies in seconds, partitioned by status code, method and route.",
		Buckets: prometheus.DefBuckets,
	}, labels)
	for _, col := range []prometheus.Collector{requests, duration} {
		if err := reg.Register(col); err != nil {
			return nil, err
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}
			// The route pattern, not the path, keeps one series per route
			// instead of one per id.
			url := r.URL.Path
			if rctx := c.RouteContext(r.Context()); rctx != nil && rctx.RoutePattern() != "" {
				url = rctx.RoutePattern()
			}
			values := []string{strconv.Itoa(status), r.Method, url}
			requests.WithLabelValues(values...).Inc()
			duration.WithLabelValues(values...).Observe(time.Since(start).Seconds())
		})
	}, nil
}
