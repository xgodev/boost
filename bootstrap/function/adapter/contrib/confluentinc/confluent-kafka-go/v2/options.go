package confluent

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

// Options can be used to create customized handler.
type Options struct {
	Topics       []string
	TimeOut      time.Duration
	ManualCommit bool
	MaxWorkers   int64
	Backoff      bool
	BackoffBase  time.Duration // Base duration for backoff
	MaxBackoff   time.Duration // Maximum backoff duration
	RetryLimit   int           // Limit for retries (-1 for infinite retries)
}

// DefaultOptions returns options based in config.
func DefaultOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
