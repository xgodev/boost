package datadog

import (
	"github.com/xgodev/boost/wrapper/config"
	awstrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go-v2/aws"
)

// Options struct which represents a datadog plugin for aws.
type Options struct {
	Enabled      bool
	TraceOptions []awstrace.Option
}

// NewOptions returns options from config file or environment vars.
func NewOptions(traceOptions ...awstrace.Option) (*Options, error) {
	opts := &Options{TraceOptions: traceOptions}
	return config.MergeOptionsWithPath[Options](opts, root)
}

// NewOptionsWithPath unmarshals options based a given key path.
func NewOptionsWithPath(path string, traceOptions ...awstrace.Option) (opts *Options, err error) {
	opts, err = NewOptions(traceOptions...)
	if err != nil {
		return nil, err
	}

	return config.MergeOptionsWithPath[Options](opts, path)
}
