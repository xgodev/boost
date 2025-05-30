package otel

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

const (
	root           = "boost.factory.otel"
	metricEnabled  = root + ".metric.enabled"
	traceEnabled   = root + ".trace.enabled"
	consoleEnabled = root + ".console.enabled"
	env            = root + ".env"
	version        = root + ".version"
	protocol       = root + ".protocol"
	endpoint       = root + ".endpoint"
	insecure       = root + ".insecure"
	export         = root + ".export"
	exportInterval = export + ".interval"
	exportTimeout  = export + ".timeout"
	attributes     = root + ".attributes"
	tlsCert        = root + ".tls.cert"
)

func init() {
	config.Add(traceEnabled, true, "enables the opentelemetry tracer")
	config.Add(metricEnabled, true, "enables the opentelemetry metrics")
	config.Add(consoleEnabled, false, "enables the opentelemetry console exporter")
	config.Add(env, "", "service env")
	config.Add(version, "0.0.0", "service version")
	config.Add(protocol, "grpc", "protocol to be used, http/grpc")
	config.Add(endpoint, "localhost:4317", `host address of the opentelemetry agent, note that this parameter will be ignored if 'OTEL_EXPORTER_OTLP_ENDPOINT' is set, 
	and the environment variable value will be used instead. See https://github.com/open-telemetry/opentelemetry-go/issues/3730`)
	config.Add(insecure, true, "enable/disable insecure connection to agent")
	config.Add(exportInterval, time.Millisecond*60000, "defines periodic reader timing for metrics")
	config.Add(exportTimeout, time.Millisecond*30000, "defines periodic reader timeout for metrics")
	config.Add(attributes, map[string]string{}, "sets a key/value pair which will be set as attributes on all spans created by the tracer. This option may be used multiple times")
	config.Add(tlsCert, "", "path to certificate to be used for tls")
}

// IsTraceEnabled returns config value from key boost.factory.otel.enabled where default is true.
func IsTraceEnabled() bool {
	return config.Bool(traceEnabled)
}

// IsMetricEnabled returns config value from key boost.factory.otel.enabled where default is true.
func IsMetricEnabled() bool {
	return config.Bool(metricEnabled)
}

func IsInsecure() bool {
	return config.Bool(insecure)
}
