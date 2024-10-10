package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/wrapper/log"
)

var (
	messagesProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "boost_function_messages_processed_total",
			Help: "Number of messages processed",
		},
		[]string{"status", "function_name"},
	)

	messageProcessingLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "boost_function_message_processing_latency_seconds",
			Help:    "Time taken to process message",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(messagesProcessed)
	prometheus.MustRegister(messageProcessingLatency)
}

type Prometheus[T any] struct {
	options *Options
}

func (c *Prometheus[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (T, error) {

	timer := prometheus.NewTimer(messageProcessingLatency)
	defer timer.ObserveDuration()

	e, err := ctx.Next(exec, fallbackFunc)
	if err != nil {
		messagesProcessed.WithLabelValues("error", c.options.FunctionName).Inc()
	} else {
		messagesProcessed.WithLabelValues("success", c.options.FunctionName).Inc()
	}

	if c.options.PushGateway.Enabled {
		if c.options.PushGateway.Async {
			go c.pushMetrics(ctx.GetContext())
		} else {
			c.pushMetrics(ctx.GetContext())
		}
	}

	return e, err
}

func (c *Prometheus[T]) pushMetrics(ctx context.Context) {
	if err := push.New(c.options.PushGateway.URL, c.options.FunctionName).
		Gatherer(prometheus.DefaultGatherer).
		Push(); err != nil {
		logger := log.FromContext(ctx).WithTypeOf(*c)
		logger.WithError(err).Warnf("error on push metrics")
	}
}

func NewAnyErrorMiddleware[T any]() (middleware.AnyErrorMiddleware[T], error) {
	return NewPrometheus[T]()
}

func NewAnyErrorMiddlewareWithOptions[T any](options *Options) middleware.AnyErrorMiddleware[T] {
	return NewPrometheusWithOptions[T](options)
}

func NewPrometheus[T any]() (*Prometheus[T], error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewPrometheusWithOptions[T](opts), nil
}

func NewPrometheusWithOptions[T any](options *Options) *Prometheus[T] {
	return &Prometheus[T]{options: options}
}
