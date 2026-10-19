package datadog

import (
	echotrace "github.com/DataDog/dd-trace-go/contrib/labstack/echo.v5/v2"
	"github.com/xgodev/boost/wrapper/config"
)

// Options datadog plugin for echo server options.
type Options struct {
	Enabled      bool
	TraceOptions []echotrace.Option
}

// NewOptions returns options from config file or environment vars.
func NewOptions(traceOptions ...echotrace.Option) (*Options, error) {
	opts := &Options{TraceOptions: traceOptions}

	return config.MergeOptionsWithPath[Options](opts, root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string, traceOptions ...echotrace.Option) (opts *Options, err error) {
	opts, err = NewOptions(traceOptions...)
	if err != nil {
		return nil, err
	}

	return config.MergeOptionsWithPath[Options](opts, path)
}
