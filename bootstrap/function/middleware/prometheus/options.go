package prometheus

import (
	"github.com/xgodev/boost/wrapper/config"
)

type Options struct {
	FunctionName string
	PushGateway  struct {
		Enabled bool
		URL     string `config:"url"`
		Async   bool
	}
}

func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
