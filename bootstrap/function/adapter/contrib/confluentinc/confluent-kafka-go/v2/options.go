package confluent

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

// Options can be used to create customized handler.
type Options struct {
	Topics  []string
	TimeOut time.Duration
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
