package datadog

import (
	"github.com/xgodev/boost"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// Options datadog plugin for grpc client options.
type Options struct {
	Enabled      bool
	traceOptions []grpctrace.Option
}

// NewOptions returns options from config file or environment vars.
func NewOptions(traceOptions ...grpctrace.Option) (*Options, error) {
	opts := &Options{
		traceOptions: traceOptions,
	}

	return boost.MergeOptionsWithPath[Options](opts, root)
}

// NewOptionsWithPath unmarshals options based a given key path.
func NewOptionsWithPath(path string, grpcOptions ...grpctrace.Option) (opts *Options, err error) {

	opts, err = NewOptions(grpcOptions...)
	if err != nil {
		return nil, err
	}

	return boost.MergeOptionsWithPath[Options](opts, path)
}
