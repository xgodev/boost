package mongo

import (
	"github.com/xgodev/boost"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Options represents mongo client options.
type Options struct {
	Uri  string
	Auth *options.Credential
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return boost.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return boost.NewOptionsWithPath[Options](root, path)
}
