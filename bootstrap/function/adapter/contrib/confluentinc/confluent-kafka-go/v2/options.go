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
	opts, err := config.NewOptionsWithPath[Options](root)
	if err != nil {
		return nil, err
	}
	// Backward compat (issue #48): merge any values still set under the legacy
	// ".kafka_confluent" root via a config file. Absent -> no-op. Deprecated.
	return config.MergeOptionsWithPath[Options](opts, legacyRoot)
}
