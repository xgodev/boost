package ignore_errors

import (
	"github.com/xgodev/boost/wrapper/config"
)

type Options struct {
	Errors []string
}

func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
