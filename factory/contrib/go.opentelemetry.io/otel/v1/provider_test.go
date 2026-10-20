package otel_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1/oteltest"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// The providers start once per process and Shutdown is final, so the whole
// life of the providers -- read before Start, Start, record, Shutdown -- is
// one test.
//
// Plugins read TracerProvider/MeterProvider while their component is built,
// which in an fx app happens before anything starts the providers. A
// batching processor and a periodic reader hold what was recorded until
// their next export, which Shutdown must not lose.
func TestProvidersForwardFromBeforeStartAndShutdownFlushes(t *testing.T) {
	collector, err := oteltest.NewCollector()
	require.NoError(t, err)
	t.Cleanup(collector.Close)
	t.Setenv("BOOST_FACTORY_OTEL_ENDPOINT", collector.Endpoint())

	boost.Start()
	ctx := context.Background()

	tracer := otelboost.TracerProvider.Tracer("boost-test")
	counter, err := otelboost.MeterProvider.Meter("boost-test").Int64Counter("boost.test.calls")
	require.NoError(t, err)

	opts, err := otelboost.NewOptions()
	require.NoError(t, err)

	spans := tracetest.NewInMemoryExporter()
	pending := keptSpans{tracetest.NewInMemoryExporter()}
	otelboost.StartTracerProviderWithOptions(ctx, opts,
		sdktrace.WithSyncer(spans),
		sdktrace.WithBatcher(pending, sdktrace.WithBatchTimeout(time.Hour)))

	reader := sdkmetric.NewManualReader()
	periodic := &metricExporter{}
	otelboost.StartMetricProviderWithOptions(ctx, opts,
		sdkmetric.WithReader(reader),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(periodic, sdkmetric.WithInterval(time.Hour))))

	_, span := tracer.Start(ctx, "boost.test.op")
	span.End()
	counter.Add(ctx, 1)

	t.Run("spans", func(t *testing.T) {
		require.Contains(t, spanNames(spans), "boost.test.op")
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

	t.Run("the batch and the periodic reader still hold them", func(t *testing.T) {
		require.Empty(t, spanNames(pending.InMemoryExporter))
		require.Empty(t, periodic.exported())
	})

	require.NoError(t, otelboost.Shutdown(ctx))

	t.Run("shutdown exports the pending spans", func(t *testing.T) {
		require.Contains(t, spanNames(pending.InMemoryExporter), "boost.test.op")
	})

	t.Run("shutdown exports the pending metrics", func(t *testing.T) {
		require.Contains(t, periodic.exported(), "boost.test.calls")
	})

	t.Run("shutdown flushes the OTLP exporters to the collector", func(t *testing.T) {
		require.Contains(t, collector.Spans(), "boost.test.op")
		require.Contains(t, collector.Metrics(), "boost.test.calls")
	})

	t.Run("nothing is exported after shutdown", func(t *testing.T) {
		_, span := tracer.Start(ctx, "boost.test.after")
		span.End()
		require.NotContains(t, spanNames(spans), "boost.test.after", "the syncer would have exported it at End")

		_, fresh := otelboost.TracerProvider.Tracer("boost-test").Start(ctx, "boost.test.fresh")
		require.False(t, fresh.IsRecording(), "a tracer taken after shutdown is a noop")
	})

	t.Run("start after shutdown does nothing", func(t *testing.T) {
		again := tracetest.NewInMemoryExporter()
		otelboost.StartTracerProviderWithOptions(ctx, opts, sdktrace.WithSyncer(again))

		_, span := otelboost.TracerProvider.Tracer("boost-test").Start(ctx, "boost.test.again")
		span.End()
		require.False(t, span.IsRecording())
		require.Empty(t, again.GetSpans())
	})

	t.Run("a second shutdown is a no-op", func(t *testing.T) {
		require.NoError(t, otelboost.Shutdown(ctx))
	})
}

func spanNames(e *tracetest.InMemoryExporter) []string {
	var names []string
	for _, s := range e.GetSpans() {
		names = append(names, s.Name)
	}
	return names
}

// keptSpans keeps the spans after Shutdown; tracetest's exporter resets
// them there, which would hide what the shutdown flushed.
type keptSpans struct{ *tracetest.InMemoryExporter }

func (keptSpans) Shutdown(context.Context) error { return nil }

// metricExporter records the names of the metrics it was asked to export.
type metricExporter struct {
	mu    sync.Mutex
	names []string
}

func (e *metricExporter) Temporality(k sdkmetric.InstrumentKind) metricdata.Temporality {
	return sdkmetric.DefaultTemporalitySelector(k)
}

func (e *metricExporter) Aggregation(k sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return sdkmetric.DefaultAggregationSelector(k)
}

func (e *metricExporter) Export(_ context.Context, rm *metricdata.ResourceMetrics) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			e.names = append(e.names, m.Name)
		}
	}
	return nil
}

func (e *metricExporter) exported() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]string(nil), e.names...)
}

func (e *metricExporter) ForceFlush(context.Context) error { return nil }

func (e *metricExporter) Shutdown(context.Context) error { return nil }
