package kafka

import (
	"github.com/xgodev/boost/config"
	"time"
)

// Options can be used to create customized handler.
type Options struct {
	Brokers          []string
	Subjects         []string
	GroupId          string
	Concurrency      int
	QueueCapacity    int
	MinBytes         int
	MaxBytes         int
	StartOffset      int64
	ReadBatchTimeout time.Duration
	MaxWait          time.Duration
}

// DefaultOptions returns options based in config.
func DefaultOptions() (*Options, error) {

	o := &Options{}

	err := config.UnmarshalWithPath(root, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}
