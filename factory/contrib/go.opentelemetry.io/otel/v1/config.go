package otel

import "github.com/xgodev/boost/wrapper/config"

const (
	root     = "boost.factory.otel"
	enabled  = root + ".enabled"
	service  = root + ".service"
	env      = root + ".env"
	version  = root + ".version"
	protocol = root + ".protocol"
	endpoint = root + ".endpoint"
	insecure = root + ".insecure"
	tags     = root + ".tags"
	tlsCert  = root + ".tls.cert"
)

func init() {
	config.Add(service, "", "service name for opentelemetry spans")
	config.Add(enabled, true, "enables the opentelemetry tracer")
	config.Add(env, "", "service env")
	config.Add(version, "0.0.0", "service version")
	config.Add(protocol, "http", "protocol to be used, http/grpc")
	config.Add(endpoint, "localhost:4318", `host address of the opentelemetry agent, note that this parameter will be ignored if 'OTEL_EXPORTER_OTLP_ENDPOINT' is set, 
	and the environment variable value will be used instead. See https://github.com/open-telemetry/opentelemetry-go/issues/3730`)
	config.Add(insecure, false, "enable/disable insecure connection to agent")
	config.Add(tags, map[string]string{}, "sets a key/value pair which will be set as a tag on all spans created by tracer. This option may be used multiple times")
	config.Add(tlsCert, "", "path to certificate to be used for tls")
}

// IsTracerEnabled returns config value from key boost.opentelemetry.enabled where default is true.
func IsTracerEnabled() bool {
	return config.Bool(enabled)
}

func IsInsecure() bool {
	return config.Bool(insecure)
}

// Service returns config value from key boost.opentelemetry.service where default is empty.
func Service() string {
	return config.String(service)
}
