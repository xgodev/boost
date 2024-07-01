package otelresty // import "github.com/xgodev/boost/factory/go-resty/resty.v2/plugins/contrib/dubonzi/otelresty.v1"

import (
	dubresty "github.com/dubonzi/otelresty"
	"github.com/xgodev/boost/wrapper/config"
)

type Options struct {
	Enabled        bool
	TracerName     string
	TracingOptions []dubresty.Option
}

// NewOptions returns options from config file or environment vars.
func NewOptions(tracingOptions ...dubresty.Option) (*Options, error) {
	opts := &Options{
		TracingOptions: tracingOptions,
	}

	return config.MergeOptionsWithPath[Options](opts, root)
}

// NewOptionsWithPath unmarshals options based a given key path.
func NewOptionsWithPath(path string, tracingOptions ...dubresty.Option) (opts *Options, err error) {
	opts, err = NewOptions(tracingOptions...)
	if err != nil {
		return nil, err
	}

	return config.MergeOptionsWithPath[Options](opts, path)
}
