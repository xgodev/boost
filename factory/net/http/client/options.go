package client

import (
	"time"

	"github.com/xgodev/boost/factory"
)

// Options http client options
type Options struct {
	MaxIdleConnPerHost    int
	MaxIdleConn           int
	MaxConnsPerHost       int
	IdleConnTimeout       time.Duration
	DisableKeepAlives     bool
	DisableCompression    bool
	ForceHTTP2            bool          `config:"forceHTTP2"`
	TLSHandshakeTimeout   time.Duration `config:"TLSHandshakeTimeout"`
	Timeout               time.Duration
	KeepAlive             time.Duration
	ExpectContinueTimeout time.Duration
	DualStack             bool
	DialTimeout           time.Duration
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return factory.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return factory.NewOptionsWithPath[Options](root, path)
}
