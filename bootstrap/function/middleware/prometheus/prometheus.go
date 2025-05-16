package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/extra/middleware"
	p "github.com/xgodev/boost/factory/contrib/prometheus/client_golang/v1"
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
	p.StartPusher()
	prometheus.MustRegister(messagesProcessed)
	prometheus.MustRegister(messageProcessingLatency)
}

type Prometheus[T any] struct {
}

func (c *Prometheus[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (T, error) {

	timer := prometheus.NewTimer(messageProcessingLatency)
	defer timer.ObserveDuration()

	e, err := ctx.Next(exec, fallbackFunc)
	var status string
	if err != nil {
		status = "error"
	} else {
		status = "success"
	}

	messagesProcessed.WithLabelValues(status, boost.ApplicationName()).Inc()

	return e, err
}

func NewAnyErrorMiddleware[T any]() middleware.AnyErrorMiddleware[T] {
	return NewPrometheus[T]()
}

func NewPrometheus[T any]() *Prometheus[T] {
	return &Prometheus[T]{}
}
