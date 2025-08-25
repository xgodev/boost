package datadog

import (
	sqltrace "github.com/DataDog/dd-trace-go/contrib/database/sql/v2"
	"github.com/xgodev/boost/wrapper/config"
)

// Options represents datadog plugin for go driver for oracle options.
type Options struct {
	Enabled      bool
	TraceOptions []sqltrace.Option
}

// NewOptions returns options from config file or environment vars.
func NewOptions(traceOptions ...sqltrace.Option) (*Options, error) {
	opts := &Options{
		TraceOptions: traceOptions,
	}

	return config.MergeOptionsWithPath[Options](opts, root)
}

// NewOptionsWithPath unmarshals options based a given key path.
func NewOptionsWithPath(path string, traceOptions ...sqltrace.Option) (opts *Options, err error) {
	opts, err = NewOptions(traceOptions...)
	if err != nil {
		return nil, err
	}

	return config.MergeOptionsWithPath[Options](opts, path)
}
