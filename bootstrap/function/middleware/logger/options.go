package logger

import (
	"github.com/xgodev/boost/wrapper/config"
)

type Options struct {
	Level      string
	ErrorStack bool
}

func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
