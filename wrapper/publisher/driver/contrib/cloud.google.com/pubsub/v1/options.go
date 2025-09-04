package pubsub

import (
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/xgodev/boost/wrapper/config"
)

// Options kafka connection options.
type Options struct {
	Log struct {
		Level string
	}
	OrderingKey    bool
	Settings       pubsub.PublishSettings
	Timeout        time.Duration
	PublishTimeout time.Duration
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
