package datadog

import (
	"github.com/xgodev/boost"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v7"
)

// Options represents the datadog client for redis cluster client options.
type Options struct {
	Enabled      bool
	TraceOptions []redistrace.ClientOption
}

// NewOptions returns options from config or environment vars.
func NewOptions(traceOptions ...redistrace.ClientOption) (*Options, error) {
	opts := &Options{TraceOptions: traceOptions}
	return boost.MergeOptionsWithPath[Options](opts, root)
}

// NewOptionsWithPath unmarshals options based a given key path.
func NewOptionsWithPath(path string, traceOptions ...redistrace.ClientOption) (opts *Options, err error) {
	opts, err = NewOptions(traceOptions...)
	if err != nil {
		return nil, err
	}

	return boost.MergeOptionsWithPath[Options](opts, path)
}
