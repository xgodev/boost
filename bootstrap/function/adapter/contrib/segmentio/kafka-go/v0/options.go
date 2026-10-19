package kafka

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

// Options of the kafka-go adapter. The reader itself (brokers, groupId,
// offsets) is configured in boost.factory.kafkago.
type Options struct {
	Topics      []string
	Buffer      int
	Backoff     bool
	BackoffBase time.Duration
	MaxBackoff  time.Duration
	RetryLimit  int
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
