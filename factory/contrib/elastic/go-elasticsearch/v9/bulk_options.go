package elasticsearch

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

// BulkOptions are the bulk indexer settings read from boost.factory.elasticsearch.bulk.
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
	return config.NewOptionsWithPath[BulkOptions](root + bulk)
}

// NewBulkOptionsWithPath returns the root bulk options overridden by the keys
// under path + ".bulk".
func NewBulkOptionsWithPath(path string) (opts *BulkOptions, err error) {
	return config.NewOptionsWithPath[BulkOptions](root+bulk, path+bulk)
}
