package otelsql

import (
	"github.com/xgodev/boost"
)

// Options represents datadog plugin for go driver for oracle options.
type Options struct {
	Enabled bool
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	opts := &Options{}

	return boost.MergeOptionsWithPath[Options](opts, root)
}

// NewOptionsWithPath unmarshals options based a given key path.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	opts, err = NewOptions()
	if err != nil {
		return nil, err
	}

	return boost.MergeOptionsWithPath[Options](opts, path)
}
