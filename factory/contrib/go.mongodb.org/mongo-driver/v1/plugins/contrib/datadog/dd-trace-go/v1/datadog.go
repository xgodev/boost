package datadog

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"

	mongotrace "github.com/DataDog/dd-trace-go/contrib/go.mongodb.org/mongo-driver.v2/v2/mongo"
	datadog "github.com/xgodev/boost/factory/contrib/datadog/dd-trace-go/v1"
	"github.com/xgodev/boost/wrapper/log"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Datadog represents a datadog plugin for mongo.
type Datadog struct {
	options *Options
}

// NewDatadogWithConfigPath returns a new datadog plugin with options from config path.
func NewDatadogWithConfigPath(path string) (*Datadog, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewDatadogWithOptions(o), nil
}

// NewDatadogWithOptions returns a new datadog plugin with options.
func NewDatadogWithOptions(options *Options) *Datadog {
	return &Datadog{options: options}
}

// NewDatadog returns a new datadog plugin with default options.
func NewDatadog() *Datadog {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewDatadogWithOptions(o)
}

// Register registers this datadog plugin on a new mongo client.
func (d *Datadog) Register(ctx context.Context) (mongo.ClientOptionsPlugin, mongo.ClientPlugin) {
	if !d.options.Enabled || !datadog.IsTracerEnabled() {
		return nil, nil
	}

	opt := []mongotrace.Option{
		mongotrace.WithService(d.options.ServiceName),
		//mongotrace.WithAnalytics(d.options.Analytics),
		//mongotrace.WithAnalyticsRate(d.options.AnalyticsRate),
	}

	return func(ctx context.Context, options *options.ClientOptions) error {
		logger := log.FromContext(ctx)

		logger.Trace("integrating datadog in mongo")

		options.SetMonitor(mongotrace.NewMonitor(opt...))

		logger.Debug("datadog successfully integrated in mongo")

		return nil
	}, nil
}

// Register registers a new datadog plugin on a new mongo client.
func Register(ctx context.Context) (mongo.ClientOptionsPlugin, mongo.ClientPlugin) {
	o, err := NewOptions()
	if err != nil {
		return nil, nil
	}
	datadog := NewDatadogWithOptions(o)
	return datadog.Register(ctx)
}
