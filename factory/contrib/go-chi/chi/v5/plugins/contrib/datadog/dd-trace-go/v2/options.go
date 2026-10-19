package datadog

import (
	chitrace "github.com/DataDog/dd-trace-go/contrib/go-chi/chi.v5/v2"
	"github.com/xgodev/boost/wrapper/config"
)

// Options datadog plugin for chi server options.
type Options struct {
	Enabled      bool
	TraceOptions []chitrace.Option
}

// NewOptions returns options from config file or environment vars.
func NewOptions(traceOptions ...chitrace.Option) (*Options, error) {
	opts := &Options{TraceOptions: traceOptions}

	return config.MergeOptionsWithPath[Options](opts, root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string, traceOptions ...chitrace.Option) (opts *Options, err error) {
	opts, err = NewOptions(traceOptions...)
	if err != nil {
		return nil, err
	}

	return config.MergeOptionsWithPath[Options](opts, path)
}
