package otel_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// Plugins read TracerProvider/MeterProvider while their component is built,
// which in an fx app happens before anything starts the providers.
func TestProvidersReadBeforeStartForwardToStartedProviders(t *testing.T) {
	boost.Start()
	ctx := context.Background()

	tracer := otelboost.TracerProvider.Tracer("boost-test")
	counter, err := otelboost.MeterProvider.Meter("boost-test").Int64Counter("boost.test.calls")
	require.NoError(t, err)

	opts, err := otelboost.NewOptions()
	require.NoError(t, err)

	spans := tracetest.NewInMemoryExporter()
	otelboost.StartTracerProviderWithOptions(ctx, opts, sdktrace.WithSyncer(spans))

	reader := sdkmetric.NewManualReader()
	otelboost.StartMetricProviderWithOptions(ctx, opts, sdkmetric.WithReader(reader))

	_, span := tracer.Start(ctx, "boost.test.op")
	span.End()
	counter.Add(ctx, 1)

	t.Run("spans", func(t *testing.T) {
		var got []string
		for _, s := range spans.GetSpans() {
			got = append(got, s.Name)
		}
		require.Contains(t, got, "boost.test.op")
	})

	t.Run("metrics", func(t *testing.T) {
		var rm metricdata.ResourceMetrics
		require.NoError(t, reader.Collect(ctx, &rm))

		var got []string
		for _, sm := range rm.ScopeMetrics {
			for _, m := range sm.Metrics {
				got = append(got, m.Name)
			}
		}
		require.Contains(t, got, "boost.test.calls")
	})
}
