# OpenTelemetry factory

Starts and stops the OpenTelemetry tracer and meter providers that every
boost OTel plugin (echo, resty, redis, mongo, grpc, pubsub, sql, mcp, ...)
reports to.

```go
import otel "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
```

## Lifecycle

| Call | What it does |
|---|---|
| `StartTracerProvider(ctx)` / `StartMeterProvider(ctx)` | Builds the SDK provider (OTLP exporter, batching span processor / periodic reader, resource) and registers it globally. Once per process; off when `boost.factory.otel.trace.enabled` / `.metric.enabled` is `false`. |
| `TracerProvider` / `MeterProvider` | What plugins instrument with. Read before Start, they forward to the provider Start registers. |
| `Shutdown(ctx)` | Exports what the providers still hold (the span batch, the metrics since the last export) and stops both, at once, within `ctx`; returns both errors joined. Final: after it nothing is exported, a Start only logs a warning, and a second Shutdown returns `nil`. |

In an fx app, include `fx/modules/factory/contrib/go.opentelemetry.io/otel/v1`:
it starts both providers and shuts them down on the app's stop, within the
stop context. A shutdown error (an unreachable collector) is logged at
ERROR and does not fail the stop, so a job that succeeded still exits 0.

Without the shutdown, a stop drops the spans still queued and the metrics
since the last interval (60s by default), and a short job can exit before
its first export.

## Config

| Key | Default | |
|---|---|---|
| `boost.factory.otel.trace.enabled` | `true` | starts the tracer provider |
| `boost.factory.otel.metric.enabled` | `true` | starts the meter provider |
| `boost.factory.otel.protocol` | `grpc` | `grpc` or `http` (`OTEL_EXPORTER_OTLP_PROTOCOL`) |
| `boost.factory.otel.endpoint` | `localhost:4317` | collector address (`OTEL_EXPORTER_OTLP_ENDPOINT`) |
| `boost.factory.otel.insecure` | `true` | plain-text connection to the collector |
| `boost.factory.otel.trace.ratio` | `1.0` | sampling ratio (`OTEL_TRACES_SAMPLER_ARG`) |
| `boost.factory.otel.export.interval` | `60s` | metric export interval (`OTEL_METRIC_EXPORT_INTERVAL`) |
| `boost.factory.otel.export.timeout` | `30s` | metric export timeout (`OTEL_METRIC_EXPORT_TIMEOUT`) |
| `boost.factory.otel.console.enabled` | `false` | also prints spans and metrics to stdout |

Per-signal endpoint and protocol (`.trace.endpoint`, `.metric.protocol`, ...)
and the resource keys (`env`, `version`, `attributes`, `tls.cert`) are in
`config.go`.

## Testing: `oteltest`

`oteltest.NewCollector()` runs an OTLP/gRPC collector in the test process
and records the names of the spans and metrics it receives. Point the
factory at it before `boost.Start()`:

```go
collector, err := oteltest.NewCollector()
require.NoError(t, err)
t.Cleanup(collector.Close)
t.Setenv("BOOST_FACTORY_OTEL_ENDPOINT", collector.Endpoint())

boost.Start()
// ... start the app, call it, stop it ...
require.Contains(t, collector.Spans(), "GET /stock/:sku")
```
