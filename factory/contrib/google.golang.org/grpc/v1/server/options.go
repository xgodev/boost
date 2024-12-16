package server

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

// Options grpc server options.
type Options struct {
	Port                  int
	MaxConcurrentStreams  int64
	InitialWindowSize     int32
	InitialConnWindowSize int32
	TLS                   TLSOptions `config:"tls"`
	KeepAlive             KeepAliveOtions
}

type KeepAliveOtions struct {
	Time                  time.Duration // Ping interval in seconds
	Timeout               time.Duration // Timeout for a ping response in seconds
	MaxConnectionIdle     time.Duration // Max idle time in seconds before closing the connection
	MaxConnectionAge      time.Duration // Max connection age in seconds before closing
	MaxConnectionAgeGrace time.Duration // Grace period in seconds before forcing connection closure
}

type TLSAutoOptions struct {
	Host string
}

type TLSOptions struct {
	Enabled bool
	Type    string
	Auto    TLSAutoOptions
	File    TLSFileOptions
}

type TLSFileOptions struct {
	Cert string
	Key  string
	CA   string `config:"ca"`
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
