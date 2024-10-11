package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
)

var (
	messagesProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "boost_wrapper_publisher_messages_sends_total",
			Help: "Number of messages sends",
		},
		[]string{"status", "source", "subject"},
	)

	messageProcessingLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "boost_wrapper_publisher_messages_sends_latency_seconds",
			Help:    "Time taken to send all messages",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(messagesProcessed)
	prometheus.MustRegister(messageProcessingLatency)
}

type Prometheus struct {
	options *Options
}

func (c *Prometheus) Exec(ctx *middleware.AnyErrorContext[[]publisher.PublishOutput], exec middleware.AnyErrorExecFunc[[]publisher.PublishOutput], fallbackFunc middleware.AnyErrorReturnFunc[[]publisher.PublishOutput]) ([]publisher.PublishOutput, error) {

	timer := prometheus.NewTimer(messageProcessingLatency)
	defer timer.ObserveDuration()

	outputs, err := ctx.Next(exec, fallbackFunc)
	if err != nil {
		return outputs, err
	}

	for _, output := range outputs {

		if output.Error != nil {
			messagesProcessed.WithLabelValues("error", output.Event.Source(), output.Event.Subject()).Inc()
		} else {
			messagesProcessed.WithLabelValues("success", output.Event.Source(), output.Event.Subject()).Inc()
		}

	}

	if c.options.PushGateway.Enabled {
		if c.options.PushGateway.Async {
			go c.pushMetrics(ctx.GetContext())
		} else {
			c.pushMetrics(ctx.GetContext())
		}
	}

	return outputs, err
}

func (c *Prometheus) pushMetrics(ctx context.Context) {
	if err := push.New(c.options.PushGateway.URL, boost.ApplicationName()).
		Gatherer(prometheus.DefaultGatherer).
		Push(); err != nil {
		logger := log.FromContext(ctx).WithTypeOf(*c)
		logger.WithError(err).Warnf("error on push metrics")
	}
}

func NewAnyErrorMiddleware() (middleware.AnyErrorMiddleware[[]publisher.PublishOutput], error) {
	return NewPrometheus()
}

func NewAnyErrorMiddlewareWithOptions(options *Options) middleware.AnyErrorMiddleware[[]publisher.PublishOutput] {
	return NewPrometheusWithOptions[[]publisher.PublishOutput](options)
}

func NewPrometheus() (*Prometheus, error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewPrometheusWithOptions[[]publisher.PublishOutput](opts), nil
}

func NewPrometheusWithOptions(options *Options) *Prometheus {
	return &Prometheus{options: options}
}
