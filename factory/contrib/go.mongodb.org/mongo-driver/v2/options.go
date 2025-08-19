package mongo

import (
	"github.com/xgodev/boost/wrapper/config"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Options represents mongo client options.
type Options struct {
	Uri  string
	Auth *options.Credential
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
