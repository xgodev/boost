package lambda

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options configures the Lambda adapter.
type Options struct {
	// Skip acknowledges every invocation without calling the handler.
	Skip bool
	// MaxWorkers bounds the records of one invocation handled at the same
	// time; 1 handles them in the order Lambda delivered them.
	MaxWorkers int
}

// NewOptions returns the options read from config.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
