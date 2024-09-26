package pubsub

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

// Options can be used to create customized handler.
type Options struct {
	Subscriptions []string
	Concurrency   int64         // Max concurrent workers
	Backoff       bool          // Enable backoff
	BackoffBase   time.Duration // Base backoff duration
	MaxBackoff    time.Duration // Max backoff duration
	RetryLimit    int           // Limit retries, -1 for infinite retries
}

// DefaultOptions returns options based in config.
func DefaultOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
