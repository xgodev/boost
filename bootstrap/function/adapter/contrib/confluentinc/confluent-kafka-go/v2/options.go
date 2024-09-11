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
}

// DefaultOptions returns options based in config.
func DefaultOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
