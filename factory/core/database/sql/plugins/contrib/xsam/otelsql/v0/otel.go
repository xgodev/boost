package otelsql

import (
	"context"
	"database/sql"
	"github.com/XSAM/otelsql"
	"github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/wrapper/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

// Register registers a new otel plugin on sql DB.
func Register(ctx context.Context, db *sql.DB) error {
	o, err := NewOptions()
	if err != nil {
		return err
	}
	h := NewOTelWithOptions(o)
	return h.Register(ctx, db)
}

// OTel represents otel plugin for go driver for oracle.
type OTel struct {
	options *Options
}

// NewOTelWithOptions returns a new otel with options.
func NewOTelWithOptions(options *Options) *OTel {
	return &OTel{options: options}
}

// NewOTelWithConfigPath returns a new otel with options from config path.
func NewOTelWithConfigPath(path string) (*OTel, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewOTelWithOptions(o), nil
}

// NewOTel returns a new otel plugin.
func NewOTel() *OTel {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewOTelWithOptions(o)
}

// Register registers this otel plugin on sql DB.
func (i *OTel) Register(ctx context.Context, db *sql.DB) error {
	if !i.options.Enabled || !otel.IsTraceEnabled() {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("integrating sql in otel")

	if err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL)); err != nil {
		return err
	}
	otelsql.WithTracerProvider(otel.TracerProvider)
	otelsql.WithMeterProvider(otel.MeterProvider)

	logger.Debug("otel successfully integrated in sql")

	return nil
}
