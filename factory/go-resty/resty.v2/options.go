package resty

import (
	"time"

	"github.com/xgodev/boost/factory"
)

// Options represents resty client options.
type Options struct {
	Debug             bool
	ConnectionTimeout time.Duration
	CloseConnection   bool
	KeepAlive         time.Duration
	RequestTimeout    time.Duration
	FallbackDelay     time.Duration
	Transport         OptionsTransport
	Host              string
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return factory.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return factory.NewOptionsWithPath[Options](root, path)
}
