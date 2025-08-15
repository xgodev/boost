package datadog

import (
	echotrace "github.com/DataDog/dd-trace-go/contrib/labstack/echo.v4/v2"
	"github.com/xgodev/boost/wrapper/config"
)

type Options struct {
	Enabled      bool
	TraceOptions []echotrace.Option
}

func NewOptions(traceOptions ...echotrace.Option) (*Options, error) {
	opts := &Options{TraceOptions: traceOptions}

	return config.MergeOptionsWithPath[Options](opts, root)
}

func NewOptionsWithPath(path string, traceOptions ...echotrace.Option) (opts *Options, err error) {

	opts, err = NewOptions(traceOptions...)
	if err != nil {
		return nil, err
	}

	return config.MergeOptionsWithPath[Options](opts, path)
}
