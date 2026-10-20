package otel

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	// tracerProvider and meterProvider are the SDK providers Start built, nil
	// when that signal is off or failed to start.
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider

	shutDown atomic.Bool
)

// Shutdown exports what the started tracer and meter providers still hold
// (the span batch, the metrics since the last export) and stops them, within
// ctx; the error joins both providers' errors. It ends OpenTelemetry for the
// process, as an SDK provider cannot be restarted: after it nothing records,
// a Start only logs a warning and a second Shutdown returns nil. A provider
// that never started is skipped.
func Shutdown(ctx context.Context) error {
	if shutDown.Swap(true) {
		return nil
	}

	// Waits for a Start in progress, and keeps a later one from building a
	// provider nothing would shut down.
	tracerOnce.Do(func() {})
	metricOnce.Do(func() {})

	// Both export at once, so an unreachable collector costs the stop one
	// export timeout, not two.
	var wg sync.WaitGroup
	var traceErr, metricErr error
	if tracerProvider != nil {
		wg.Go(func() {
			if err := tracerProvider.Shutdown(ctx); err != nil {
				traceErr = fmt.Errorf("shutting down the opentelemetry tracer provider: %w", err)
			}
		})
	}
	if meterProvider != nil {
		wg.Go(func() {
			if err := meterProvider.Shutdown(ctx); err != nil {
				metricErr = fmt.Errorf("shutting down the opentelemetry meter provider: %w", err)
			}
		})
	}
	wg.Wait()
	return errors.Join(traceErr, metricErr)
}
