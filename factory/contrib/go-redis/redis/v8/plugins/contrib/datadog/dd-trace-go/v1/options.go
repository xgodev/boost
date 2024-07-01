package datadog

import (
	"github.com/xgodev/boost/wrapper/config"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

// Options represents a datadog client for redis options.
type Options struct {
	Enabled      bool
	TraceOptions []redistrace.ClientOption
}

// NewOptions returns options from config or environment vars.
func NewOptions(traceOptions ...redistrace.ClientOption) (*Options, error) {
	opts := &Options{TraceOptions: traceOptions}
	return config.MergeOptionsWithPath[Options](opts, root)
}

// NewOptionsWithPath unmarshals options based a given key path.
func NewOptionsWithPath(path string, traceOptions ...redistrace.ClientOption) (opts *Options, err error) {
	opts, err = NewOptions(traceOptions...)
	if err != nil {
		return nil, err
	}

	return config.MergeOptionsWithPath[Options](opts, path)
}
