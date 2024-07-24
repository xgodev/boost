package otel

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc/credentials"
	"sync"
)

var MeterProvider metric.MeterProvider

// StartMeterProvider starts the tracer provider like StartMetricProviderWithOptions but with default Options.
func StartMeterProvider(ctx context.Context, startOptions ...sdkmetric.Option) {

	o, err := NewOptions()
	if err != nil {
		panic(err)
	}

	StartMetricProviderWithOptions(ctx, o, startOptions...)
}

var metricOnce sync.Once

// StartMetricProviderWithOptions starts the tracer provider with the given set of options. Calling
// it multiple times will have no effect. If an error occours during tracer initialization,
// a Noop trace provider will be used instead.
func StartMetricProviderWithOptions(ctx context.Context, options *Options, startOptions ...sdkmetric.Option) {

	if !IsMetricEnabled() {
		return
	}

	metricOnce.Do(func() {

		MeterProvider = noop.NewMeterProvider()

		logger := log.FromContext(ctx)

		otel.SetLogger(logr.New(&Logger{}))

		exporter, err := NewMeterExporter(ctx, options)
		if err != nil {
			logger.WithError(err).Errorf("error creating opentelemetry exporter")
			return
		}

		rs, err := NewResource(ctx, options)
		if err != nil {
			logger.WithError(err).Errorf("error creating opentelemetry resource")
			return
		}

		periodicReader, err := NewReader(options, exporter)
		if err != nil {
			logger.WithError(err).Errorf("error creating opentelemetry reader")
			return
		}

		startOptions = append(startOptions,
			sdkmetric.WithReader(periodicReader),
			sdkmetric.WithResource(rs),
		)

		prov := sdkmetric.NewMeterProvider(startOptions...)

		otel.SetMeterProvider(prov)
		MeterProvider = prov

		log.Infof("started opentelemetry meter: %s", options.Service)
	})
}

func NewReader(options *Options, exporter sdkmetric.Exporter) (sdkmetric.Reader, error) {

	periodicReaderOpts := []sdkmetric.PeriodicReaderOption{
		sdkmetric.WithInterval(options.Export.Interval),
		sdkmetric.WithTimeout(options.Export.Timeout),
	}

	return sdkmetric.NewPeriodicReader(exporter, periodicReaderOpts...), nil
}

func NewMeterExporter(ctx context.Context, options *Options) (sdkmetric.Exporter, error) {
	var exporter sdkmetric.Exporter
	var err error

	switch options.Protocol {
	case "grpc":
		exporter, err = NewGRPCMeterExporter(ctx, options)
	default:
		exporter, err = NewHTTPMeterExporter(ctx, options)
	}
	return exporter, err
}

func NewHTTPMeterExporter(ctx context.Context, options *Options) (sdkmetric.Exporter, error) {
	exporterOpts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(options.Endpoint),
	}

	if IsInsecure() {
		exporterOpts = append(exporterOpts, otlpmetrichttp.WithInsecure())
	}

	exporter, err := otlpmetrichttp.New(
		ctx,
		exporterOpts...,
	)
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func NewGRPCMeterExporter(ctx context.Context, options *Options) (sdkmetric.Exporter, error) {
	exporterOpts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(options.Endpoint),
	}

	if IsInsecure() {
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithInsecure())
	} else {
		creds, err := credentials.NewClientTLSFromFile(options.TLS.Cert, "")
		if err != nil {
			return nil, errors.Wrap(err, "error creating tls credentials")
		}
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithTLSCredentials(creds))
	}

	exporter, err := otlpmetricgrpc.New(
		ctx,
		exporterOpts...,
	)
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

// NewMeter creates a Metric with the provided name and options. A Meter
// allows for the custom instrumentation.
//
// StartMeterProvider should be called before to setup the meter provider, otherwise a Noop
// tracer provider will be used.
func NewMeter(name string, options ...metric.MeterOption) metric.Meter {
	return MeterProvider.Meter(name, options...)
}
