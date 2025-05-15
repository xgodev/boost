package otelsql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/XSAM/otelsql"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/wrapper/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

// wrappedConnector delegates Connect to the original connector
// but returns the wrappedDriver from Driver().
type wrappedConnector struct {
	orig          driver.Connector
	wrappedDriver driver.Driver
}

func (w *wrappedConnector) Connect(ctx context.Context) (driver.Conn, error) {
	logger := log.FromContext(ctx)
	logger.Tracef("OTel-Connector CONNECT called, driver now: %T", w.wrappedDriver)
	return w.orig.Connect(ctx)
}

func (w *wrappedConnector) Driver() driver.Driver {
	log.Tracef("OTel-Connector DRIVER called, returning: %T", w.wrappedDriver)
	return w.wrappedDriver
}

// OTel instruments database/sql with OpenTelemetry: it wraps the connector
// to emit spans for each query and registers pool metrics on *sql.DB*.
type OTel struct {
	options *Options
}

// NewOTelWithOptions constructs the plugin with explicit Options.
func NewOTelWithOptions(options *Options) *OTel {
	return &OTel{options: options}
}

// NewOTelWithConfigPath loads Options from the given config path.
func NewOTelWithConfigPath(path string) (*OTel, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewOTelWithOptions(opts), nil
}

// NewOTel constructs the plugin with default Options.
func NewOTel() *OTel {
	opts, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return NewOTelWithOptions(opts)
}

// WrapConnector is called before sql.OpenDB. If tracing is enabled,
// it wraps the connector's Driver() so that each query emits an OTel span.
func (p *OTel) WrapConnector(ctx context.Context, connector driver.Connector) (driver.Connector, error) {
	if !p.options.Enabled || (!otelboost.IsTraceEnabled() && !otelboost.IsMetricEnabled()) {
		return connector, nil
	}
	logger := log.FromContext(ctx)
	logger.Trace("wrapping connector for OpenTelemetry SQL tracing")

	var opts []otelsql.Option
	opts = append(opts, otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL))

	if otelboost.IsTraceEnabled() {
		opts = append(opts, otelsql.WithTracerProvider(otelboost.TracerProvider))
	}

	if otelboost.IsMetricEnabled() {
		opts = append(opts, otelsql.WithMeterProvider(otelboost.MeterProvider))
	}

	origDriver := connector.Driver()
	wrappedDriver := otelsql.WrapDriver(
		origDriver,
		opts...,
	)

	return &wrappedConnector{
		orig:          connector,
		wrappedDriver: wrappedDriver,
	}, nil
}

// InitDB is called immediately after sql.OpenDB. If metrics are enabled,
// it registers DB pool stats with OpenTelemetry.
func (p *OTel) InitDB(ctx context.Context, db *sql.DB) error {
	if !p.options.Enabled || !otelboost.IsMetricEnabled() {
		return nil
	}
	logger := log.FromContext(ctx)
	logger.Trace("registering OpenTelemetry SQL pool metrics")

	var opts []otelsql.Option
	opts = append(opts, otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL))

	if otelboost.IsMetricEnabled() {
		opts = append(opts, otelsql.WithMeterProvider(otelboost.MeterProvider))
	}

	return otelsql.RegisterDBStatsMetrics(
		db,
		opts...,
	)
}
