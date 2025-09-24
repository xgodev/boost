package pubsub

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options can be used to create customized handler.
type Options struct {
	Subscriptions []string
	Concurrency   int64 // Max concurrent workers
}

// DefaultOptions returns options based in config.
func DefaultOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
