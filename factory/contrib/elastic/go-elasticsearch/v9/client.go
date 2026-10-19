package elasticsearch

import (
	"context"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Plugin defines a function to process plugin.
type Plugin = factory.Plugin[*elasticsearch.Client]

// NewClient returns elasticsearch client with default options.
func NewClient(ctx context.Context, plugins ...Plugin) (*elasticsearch.Client, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithConfigPath returns elasticsearch client with options from config path.
func NewClientWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*elasticsearch.Client, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, opts, plugins...)
}

// NewClientWithOptions builds the client, pings the cluster and runs the
// plugins in order on it. A ping the cluster answers with an error status is
// a ServiceUnavailable; a plugin error is returned and the later plugins do
// not run.
func NewClientWithOptions(ctx context.Context, o *Options, plugins ...Plugin) (*elasticsearch.Client, error) {
	client, err := elasticsearch.NewClient(clientConfig(o))
	if err != nil {
		return nil, err
	}

	res, err := client.Ping(client.Ping.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, errors.ServiceUnavailablef("elasticsearch %s answered the ping with %s", strings.Join(o.Addresses, ","), res.Status())
	}

	if err := factory.ApplyPlugins(ctx, client, plugins...); err != nil {
		return nil, err
	}

	log.FromContext(ctx).Infof("Connected to Elastic Search server: %v status: %s", strings.Join(o.Addresses, ","), res.Status())

	return client, nil
}

// clientConfig maps the options onto the go-elasticsearch config.
func clientConfig(o *Options) elasticsearch.Config {
	cfg := elasticsearch.Config{
		Addresses:             o.Addresses,
		Username:              o.Username,
		Password:              o.Password,
		CloudID:               o.CloudID,
		APIKey:                o.APIKey,
		RetryOnStatus:         o.RetryOnStatus,
		DisableRetry:          o.DisableRetry,
		MaxRetries:            o.MaxRetries,
		DiscoverNodesOnStart:  o.DiscoverNodesOnStart,
		DiscoverNodesInterval: o.DiscoverNodesInterval,
		EnableMetrics:         o.EnableMetrics,
		EnableDebugLogger:     o.EnableDebugLogger,
		RetryBackoff:          backOff(o.RetryBackoff),
		Logger:                &Logger{},
	}

	if o.CACert != "" {
		cfg.CACert = []byte(o.CACert)
	}

	return cfg
}

// backOff waits step times the attempt number before each retry.
func backOff(step time.Duration) func(attempt int) time.Duration {
	return func(attempt int) time.Duration {
		return time.Duration(attempt) * step
	}
}
