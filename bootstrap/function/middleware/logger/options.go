package logger

import (
	"github.com/xgodev/boost/wrapper/config"
)

type Options struct {
	Level      string
	ErrorStack bool
}

func NewOptions() (*Options, error) {
	o := &Options{}
	err := config.UnmarshalWithPath(root, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
