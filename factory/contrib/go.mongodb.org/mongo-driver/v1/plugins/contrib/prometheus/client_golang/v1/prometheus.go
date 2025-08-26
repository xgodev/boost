package prometheus

import (
	"context"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	"github.com/xgodev/boost/wrapper/log"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	labelNames = []string{"command"}

	cmds = prom.NewHistogramVec(prom.HistogramOpts{
		Name:    "boost_factory_mongo_commands",
		Help:    "Histogram of MongoDB commands",
		Buckets: prom.DefBuckets,
	}, labelNames)

	cmderr = prom.NewCounterVec(prom.CounterOpts{
		Name: "boost_factory_mongo_command_errors",
		Help: "Number of MongoDB commands that have failed",
	}, labelNames)
)

func init() {
	prom.MustRegister(cmderr)
	prom.MustRegister(cmds)
}

// Prometheus represents a prometheus plugin for mongo.
type Prometheus struct {
	options *Options
}

// NewPrometheusWithConfigPath returns a new prometheus plugin with options from config path.
func NewPrometheusWithConfigPath(path string) (*Prometheus, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewPrometheusWithOptions(o), nil
}

// NewPrometheusWithOptions returns a new prometheus plugin with options.
func NewPrometheusWithOptions(options *Options) *Prometheus {
	return &Prometheus{options: options}
}

// NewPrometheus returns a new prometheus plugin with default options.
func NewPrometheus() *Prometheus {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewPrometheusWithOptions(o)
}

// Register registers this prometheus plugin on a new mongo client.
func (d *Prometheus) Register(ctx context.Context) (mongo.ClientOptionsPlugin, mongo.ClientPlugin) {
	if !d.options.Enabled {
		return nil, nil
	}

	return func(ctx context.Context, options *options.ClientOptions) error {
		logger := log.FromContext(ctx)

		logger.Trace("integrating prometheus in mongo")

		options.SetMonitor(monitor())

		logger.Debug("prometheus successfully integrated in mongo")

		return nil
	}, nil
}

// Register registers a new prometheus plugin on a new mongo client.
func Register(ctx context.Context) (mongo.ClientOptionsPlugin, mongo.ClientPlugin) {
	o, err := NewOptions()
	if err != nil {
		return nil, nil
	}
	plugin := NewPrometheusWithOptions(o)
	return plugin.Register(ctx)
}

func monitor() *event.CommandMonitor {

	observeDuration := func(evt event.CommandFinishedEvent) {
		duration := evt.Duration.Seconds()
		cmds.WithLabelValues(evt.CommandName).Observe(duration)
	}

	return &event.CommandMonitor{
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			observeDuration(evt.CommandFinishedEvent)
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			observeDuration(evt.CommandFinishedEvent)
			cmderr.WithLabelValues(evt.CommandName).Inc()
		},
	}
}
