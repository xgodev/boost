package nats

import "github.com/xgodev/boost/wrapper/config"

// Options can be used to create customized handler.
type Options struct {
	Subjects []string
	Queue    string
}

// NewOptions returns options based in config.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
