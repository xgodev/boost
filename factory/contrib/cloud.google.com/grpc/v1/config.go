package grpc

import "github.com/xgodev/boost/wrapper/config"

// ConfigAdd registers shared gRPC client configuration keys under the given root path.
func ConfigAdd(path string) {
	config.Add(path+".host", "", "gRPC server host")
	config.Add(path+".port", "0", "gRPC server port")
	config.Add(path+".tls.enabled", "false", "enable TLS")
	config.Add(path+".tls.certFile", "", "path to TLS certificate file")
	config.Add(path+".tls.keyFile", "", "path to TLS key file")
	config.Add(path+".tls.caFile", "", "path to CA certificate file")
	config.Add(path+".tls.insecureSkipVerify", "false", "skip TLS certificate verification")
	config.Add(path+".initialWindowSize", "65535", "initial stream window size in bytes")
	config.Add(path+".initialConnWindowSize", "1048576", "initial connection window size in bytes")
	config.Add(path+".block", "false", "block until the connection is ready")
	config.Add(path+".hostOverwrite", "", "authority header override")
	config.Add(path+".connectParams.backoff.baseDelay", "1s", "gRPC backoff base delay")
	config.Add(path+".connectParams.backoff.multiplier", "1.6", "gRPC backoff multiplier")
	config.Add(path+".connectParams.backoff.jitter", "0.2", "gRPC backoff jitter")
	config.Add(path+".connectParams.backoff.maxDelay", "120s", "gRPC backoff maximum delay")
	config.Add(path+".connectParams.minConnectTimeout", "20s", "minimum connection timeout")
	config.Add(path+".keepalive.time", "30s", "gRPC keepalive ping interval")
	config.Add(path+".keepalive.timeout", "10s", "gRPC keepalive ping timeout")
	config.Add(path+".keepalive.permitWithoutStream", "false", "permit keepalive pings without active streams")
}
