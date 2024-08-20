package ignore_errors

import (
	"github.com/xgodev/boost/wrapper/config"
)

type Options struct {
	Errors []string
}

func NewOptions() (*Options, error) {
	o := &Options{}
	err := config.UnmarshalWithPath(root, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
