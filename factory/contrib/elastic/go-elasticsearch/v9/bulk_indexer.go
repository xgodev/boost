package elasticsearch

import (
	"context"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esutil"
	"github.com/xgodev/boost/wrapper/log"
)

// NewBulkIndexer returns elasticsearch bulk indexer with default options.
func NewBulkIndexer(ctx context.Context, client *elasticsearch.Client) (esutil.BulkIndexer, error) {
	o, err := NewBulkOptions()
	if err != nil {
		return nil, err
	}
	return NewBulkIndexerWithOptions(ctx, o, client)
}

// NewBulkIndexerWithConfigPath returns elasticsearch bulk indexer with options from config path.
func NewBulkIndexerWithConfigPath(ctx context.Context, path string, client *elasticsearch.Client) (esutil.BulkIndexer, error) {
	opts, err := NewBulkOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewBulkIndexerWithOptions(ctx, opts, client)
}

// NewBulkIndexerWithOptions returns a bulk indexer on client; a failed flush
// is logged as an error on the ctx logger.
func NewBulkIndexerWithOptions(ctx context.Context, options *BulkOptions, client *elasticsearch.Client) (esutil.BulkIndexer, error) {
	logger := log.FromContext(ctx)

	logger.Infof("Creating Elastic Bulk Indexer on index %s", options.Index)

	return esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		NumWorkers:    options.NumWorkers,
		FlushBytes:    options.FlushBytes,
		FlushInterval: options.FlushInterval,
		Client:        client,
		DebugLogger:   &DebugLogger{},
		OnError: func(ctx context.Context, err error) {
			logger.Errorf("ERROR: %s", err.Error())
		},
		Index:      options.Index,
		ErrorTrace: true,
		Pipeline:   options.Pipeline,
		Timeout:    options.Timeout,
	})
}
