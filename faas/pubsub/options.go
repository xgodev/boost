package pubsub

import (
	"github.com/xgodev/boost/config"
)

// Options can be used to create customized handler.
type Options struct {
	Topic       string
	Concurrency int
}

// DefaultOptions returns options based in config.
func DefaultOptions() (*Options, error) {

	o := &Options{}

	err := config.UnmarshalWithPath(root, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}
