package otel

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
)

type Options struct {
	Enabled bool
	Service string
	Env     string
	Version string
	Export  struct {
		Interval time.Duration
		Timeout  time.Duration
	}
	Console struct {
		Enabled bool
	}
	Protocol   string
	Endpoint   string
	Insecure   bool
	Metric     struct {
		Endpoint string
		Protocol string
	}
	Trace struct {
		Endpoint string
		Protocol string
		Ratio    float64
	}
	Attributes map[string]string
	TLS        struct {
		Cert string
	}
}

// NewOptionsWithPath unmarshals options based on a given key path.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	opts, err = NewOptions()
	if err != nil {
		return nil, err
	}

	err = config.UnmarshalWithPath(path, opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

// NewOptions returns options from config file or environment vars
// following the official OpenTelemetry environment variable specification.
// Resolution order: signal-specific env → general env → config → default.
func NewOptions() (*Options, error) {
	opts, err := config.NewOptionsWithPath[Options](root)
	if err != nil {
		return nil, err
	}

	opts.Service = boost.ApplicationName()

	// General endpoint: OTEL_EXPORTER_OTLP_ENDPOINT → config → default
	if v := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); v != "" {
		opts.Endpoint = v
	}

	// General protocol: OTEL_EXPORTER_OTLP_PROTOCOL → config → default
	if v := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"); v != "" {
		opts.Protocol = v
	}

	// Service name: OTEL_SERVICE_NAME → config → default (boost.ApplicationName)
	if v := os.Getenv("OTEL_SERVICE_NAME"); v != "" {
		opts.Service = v
	}

	// Version: OTEL_SERVICE_VERSION → config → default
	if v := os.Getenv("OTEL_SERVICE_VERSION"); v != "" {
		opts.Version = v
	}

	// Insecure: OTEL_EXPORTER_OTLP_INSECURE → config → default
	if v := os.Getenv("OTEL_EXPORTER_OTLP_INSECURE"); v != "" {
		opts.Insecure = v == "true"
	}

	// Environment/deployment from OTEL_RESOURCE_ATTRIBUTES
	// Standard: OTEL_RESOURCE_ATTRIBUTES=deployment.environment=production
	if v := os.Getenv("OTEL_RESOURCE_ATTRIBUTES"); v != "" {
		for _, kv := range strings.Split(v, ",") {
			kv = strings.TrimSpace(kv)
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) == 2 && parts[0] == "deployment.environment" {
				opts.Env = parts[1]
				break
			}
		}
	}

	// Metric-specific endpoint: OTEL_EXPORTER_OTLP_METRICS_ENDPOINT → general → config → default
	if v := os.Getenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"); v != "" {
		opts.Metric.Endpoint = v
	} else {
		opts.Metric.Endpoint = opts.Endpoint
	}

	// Metric-specific protocol: OTEL_EXPORTER_OTLP_METRICS_PROTOCOL → general → config → default
	if v := os.Getenv("OTEL_EXPORTER_OTLP_METRICS_PROTOCOL"); v != "" {
		opts.Metric.Protocol = v
	} else {
		opts.Metric.Protocol = opts.Protocol
	}

	// Trace-specific endpoint: OTEL_EXPORTER_OTLP_TRACES_ENDPOINT → general → config → default
	if v := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"); v != "" {
		opts.Trace.Endpoint = v
	} else {
		opts.Trace.Endpoint = opts.Endpoint
	}

	// Trace-specific protocol: OTEL_EXPORTER_OTLP_TRACES_PROTOCOL → general → config → default
	if v := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"); v != "" {
		opts.Trace.Protocol = v
	} else {
		opts.Trace.Protocol = opts.Protocol
	}

	// Trace sampling ratio: OTEL_TRACES_SAMPLER_ARG → config → default (1.0)
	if v := os.Getenv("OTEL_TRACES_SAMPLER_ARG"); v != "" {
		if r, err := fmt.Sscanf(v, "%f", &opts.Trace.Ratio); err != nil || r != 1 {
			return nil, fmt.Errorf("OTEL_TRACES_SAMPLER_ARG: invalid ratio %q", v)
		}
	}

	// Export interval: OTEL_METRIC_EXPORT_INTERVAL → config → default
	if v := os.Getenv("OTEL_METRIC_EXPORT_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("OTEL_METRIC_EXPORT_INTERVAL: %w", err)
		}
		opts.Export.Interval = d
	}

	// Export timeout: OTEL_METRIC_EXPORT_TIMEOUT → config → default
	if v := os.Getenv("OTEL_METRIC_EXPORT_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("OTEL_METRIC_EXPORT_TIMEOUT: %w", err)
		}
		opts.Export.Timeout = d
	}

	// OTEL_EXPORTER_OTLP_TIMEOUT as fallback for export timeout
	if v := os.Getenv("OTEL_EXPORTER_OTLP_TIMEOUT"); v != "" {
		if v2 := os.Getenv("OTEL_METRIC_EXPORT_TIMEOUT"); v2 == "" {
			d, err := time.ParseDuration(v)
			if err != nil {
				return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_TIMEOUT: %w", err)
			}
			opts.Export.Timeout = d
		}
	}

	return opts, nil
}
