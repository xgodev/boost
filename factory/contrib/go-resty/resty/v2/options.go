package resty

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

// Options represents resty client options.
type Options struct {
	Debug             bool
	Accept            string
	Authorization     string
	Headers           map[string]string
	QueryParams       map[string]string
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
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
