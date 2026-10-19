package kafka

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

// Options kafka-go options.
type Options struct {
	Brokers []string
	Writer  WriterOptions
	Reader  ReaderOptions
	Conn    ConnOptions
}

// WriterOptions options of a kafka.Writer.
type WriterOptions struct {
	Balancer               string
	RequiredAcks           string
	Compression            string
	MaxAttempts            int
	BatchSize              int
	BatchBytes             int64
	BatchTimeout           time.Duration
	ReadTimeout            time.Duration
	WriteTimeout           time.Duration
	Async                  bool
	AllowAutoTopicCreation bool
}

// ReaderOptions options of a kafka.Reader.
type ReaderOptions struct {
	GroupID          string
	Topic            string
	Partition        int
	QueueCapacity    int
	MinBytes         int
	MaxBytes         int
	MaxWait          time.Duration
	ReadBatchTimeout time.Duration
	CommitInterval   time.Duration
	StartOffset      string
}

// ConnOptions options of a kafka.Conn.
type ConnOptions struct {
	Type      string
	Network   string
	Topic     string
	Partition int
	Timeout   time.Duration
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (*Options, error) {
	return config.NewOptionsWithPath[Options](root, path)
}
