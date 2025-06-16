package datadog

import (
	"context"
	"database/sql"
	"database/sql/driver"

	ddboost "github.com/xgodev/boost/factory/contrib/datadog/dd-trace-go/v1"
	"github.com/xgodev/boost/wrapper/log"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

// Datadog instruments database/sql for DataDog tracing.
type Datadog struct {
	options *Options
}

// NewDatadogWithOptions constructs the plugin with explicit Options.
func NewDatadogWithOptions(options *Options) *Datadog {
	return &Datadog{options: options}
}

// NewDatadogWithConfigPath loads Options (and traceOpts) from the given path.
func NewDatadogWithConfigPath(path string, traceOptions ...sqltrace.Option) (*Datadog, error) {
	opts, err := NewOptionsWithPath(path, traceOptions...)
	if err != nil {
		return nil, err
	}
	return NewDatadogWithOptions(opts), nil
}

// NewDatadog constructs the plugin with default Options.
func NewDatadog(traceOptions ...sqltrace.Option) *Datadog {
	opts, err := NewOptions(traceOptions...)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return NewDatadogWithOptions(opts)
}

// WrapConnector is a no-op for DataDog: instrumentation happens on InitDB.
func (d *Datadog) WrapConnector(ctx context.Context, connector driver.Connector) (driver.Connector, error) {
	return connector, nil
}

// InitDB registers the DataDog SQL driver to start tracing queries.
func (d *Datadog) InitDB(ctx context.Context, db *sql.DB) error {
	if !d.options.Enabled || !ddboost.IsTracerEnabled() {
		return nil
	}
	logger := log.FromContext(ctx)
	logger.Trace("integrating sql in datadog")

	// Register the existing driver under the given name for tracing
	sqltrace.Register(
		"datadog-sql",
		db.Driver(),
		d.options.TraceOptions...,
	)

	logger.Debug("datadog successfully integrated in sql")
	return nil
}
