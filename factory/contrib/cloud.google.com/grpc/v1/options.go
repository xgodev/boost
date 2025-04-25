package grpc

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

// Options holds shared gRPC client configuration.
type Options struct {
	Host string `config:"host"`
	Port int    `config:"port"`
	TLS  struct {
		Enabled            bool   `config:"enabled"`
		CertFile           string `config:"certFile"`
		KeyFile            string `config:"keyFile"`
		CAFile             string `config:"caFile"`
		InsecureSkipVerify bool   `config:"insecureSkipVerify"`
	} `config:"tls"`
	InitialWindowSize     int    `config:"initialWindowSize"`
	InitialConnWindowSize int    `config:"initialConnWindowSize"`
	Block                 bool   `config:"block"`
	HostOverwrite         string `config:"hostOverwrite"`
	Backoff               struct {
		BaseDelay  time.Duration `config:"baseDelay"`
		Multiplier float64       `config:"multiplier"`
		Jitter     float64       `config:"jitter"`
		MaxDelay   time.Duration `config:"maxDelay"`
	} `config:"backoff"`
	MinConnectTimeout time.Duration `config:"minConnectTimeout"`
	Keepalive         struct {
		Time                time.Duration `config:"time"`
		Timeout             time.Duration `config:"timeout"`
		PermitWithoutStream bool          `config:"permitWithoutStream"`
	} `config:"keepalive"`
}

// NewOptionsWithPath loads Options from the specified config root.
func NewOptionsWithPath(path string) (*Options, error) {
	return config.NewOptionsWithPath[Options](path)
}
