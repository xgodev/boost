package otel

import (
	"github.com/xgodev/boost"
	"os"
	"time"

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
	Protocol   string
	Endpoint   string
	Insecure   bool
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

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {

	opts, err := config.NewOptionsWithPath[Options](root)
	if err != nil {
		return nil, err
	}

	if v := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); v != "" {
		opts.Endpoint = v
	} else {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", opts.Endpoint)
	}

	if v := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"); v != "" {
		opts.Protocol = v
	} else {
		os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", opts.Protocol)
	}

	opts.Service = boost.ApplicationName()

	if v := os.Getenv("OTEL_SERVICE_NAME"); v != "" {
		opts.Service = v
	} else {
		os.Setenv("OTEL_SERVICE_NAME", opts.Service)
	}

	if v := os.Getenv("OTEL_SERVICE_VERSION"); v != "" {
		opts.Version = v
	}

	if v := os.Getenv("OTEL_ENV"); v != "" {
		opts.Env = v
	}

	if v := os.Getenv("OTEL_METRIC_EXPORT_INTERVAL"); v != "" {
		exportInterval, err := time.ParseDuration(v)
		if err != nil {
			return nil, err
		}

		opts.Export.Interval = exportInterval
	}

	if v := os.Getenv("OTEL_METRIC_EXPORT_TIMEOUT"); v != "" {
		exportTimeout, err := time.ParseDuration(v)
		if err != nil {
			return nil, err
		}

		opts.Export.Timeout = exportTimeout
	}

	return opts, nil
}
