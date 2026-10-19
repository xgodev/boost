package azblob

import "github.com/xgodev/boost/wrapper/config"

// Options azblob client options.
type Options struct {
	AccountName      string `config:"accountName"`
	URL              string `config:"url"`
	ConnectionString string `config:"connectionString"`
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	ConfigAdd(path)
	config.Load()
	return config.NewOptionsWithPath[Options](root, path)
}
