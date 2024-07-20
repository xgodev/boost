package otel

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc/credentials"
)

var meterProvider metric.MeterProvider

// StartMetricProvider starts the tracer provider like StartMetricProviderWithOptions but with default Options.
func StartMetricProvider(ctx context.Context, startOptions ...sdkmetric.Option) {

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

		meterProvider = noop.NewMeterProvider()

		logger := log.FromContext(ctx)

		var exporter sdkmetric.Exporter
		var err error

		switch options.Protocol {
		case "grpc":
			exporter, err = startGRPCMeter(ctx, options)
		case "http":
			exporter, err = startHTTPMeter(ctx, options)
		default:
			exporter, err = startHTTPMeter(ctx, options)
		}

		if err != nil {
			logger.Error("error creating opentelemetry exporter: ", err)
			otel.SetMeterProvider(noop.NewMeterProvider())
			return
		}

		rs, err := NewResource(ctx, options)
		if err != nil {
			logger.Error("error creating opentelemetry resource: ", err)
			otel.SetMeterProvider(noop.NewMeterProvider())
			return
		}

		startOptions = append(startOptions,
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
			sdkmetric.WithResource(rs),
		)

		prov := sdkmetric.NewMeterProvider(startOptions...)

		otel.SetMeterProvider(prov)
		meterProvider = prov

		log.Infof("started opentelemetry meter: %s", options.Service)
	})
}

func startHTTPMeter(ctx context.Context, options *Options) (sdkmetric.Exporter, error) {
	var exporterOpts []otlpmetrichttp.Option
	if _, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT"); !ok { // Only using WithEndpoint when the environment variable is not set
		exporterOpts = append(exporterOpts, otlpmetrichttp.WithEndpoint(options.Endpoint)) //TODO see https://github.com/open-telemetry/opentelemetry-go/issues/3730
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

func startGRPCMeter(ctx context.Context, options *Options) (sdkmetric.Exporter, error) {
	var exporterOpts []otlpmetricgrpc.Option
	if _, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT"); !ok { // Only using WithEndpoint when the environment variable is not set
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithEndpoint(options.Endpoint)) //TODO see https://github.com/open-telemetry/opentelemetry-go/issues/3730
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

// NewMeter creates a Tracer with the provided name and options. A Tracer
// allows the creation of spans for custom instrumentation.
//
// StartMetricProvider should be called before to setup the meter provider, otherwise a Noop
// tracer provider will be used.
func NewMeter(name string, options ...metric.MeterOption) metric.Meter {
	return meterProvider.Meter(name, options...)
}
