package lambda

import (
	"github.com/xgodev/boost/config"
)

// Options can be used to create customized handler.
type Options struct {
	Skip bool
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
