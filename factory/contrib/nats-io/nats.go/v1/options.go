package nats

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

// Options nats connection options.
type Options struct {
	Url                  string
	MaxReconnects        int
	ReconnectWait        time.Duration
	ReconnectJitter      time.Duration
	ReconnectJitterTLS   time.Duration
	Timeout              time.Duration
	PingInterval         time.Duration
	MaxPingOut           int
	MaxChanLen           int
	ReconnectBufSize     int
	DrainTimeout         time.Duration
	Verbose              bool
	Compression          bool
	RetryOnFailedConnect bool
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
