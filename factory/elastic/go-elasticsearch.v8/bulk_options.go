package elasticsearch

import (
	"time"

	"github.com/xgodev/boost/factory"
)

type BulkOptions struct {
	NumWorkers    int
	FlushInterval time.Duration
	FlushBytes    int
	Pipeline      string
	Timeout       time.Duration
	Index         string
}

// NewBulkOptions returns options from config file or environment vars.
func NewBulkOptions() (*BulkOptions, error) {
	return factory.NewOptionsWithPath[BulkOptions](root + bulk)
}

// NewBulkOptionsWithPath unmarshals a given key path into options and returns it.
func NewBulkOptionsWithPath(path string) (opts *BulkOptions, err error) {
	return factory.NewOptionsWithPath[BulkOptions](root, path+bulk)
}
