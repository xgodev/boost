package datadog

import (
	redistrace "github.com/DataDog/dd-trace-go/contrib/redis/go-redis.v9/v2"
	"github.com/xgodev/boost/wrapper/config"
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
