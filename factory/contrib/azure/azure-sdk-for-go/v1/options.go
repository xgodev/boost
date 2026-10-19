package azure

import "github.com/xgodev/boost/wrapper/config"

// Options azure credential options.
type Options struct {
	ClientId     string `config:"clientId"`
	ClientSecret string `config:"clientSecret"`
	TenantId     string `config:"tenantId"`
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](Root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	ConfigAdd(path)
	config.Load()
	return config.NewOptionsWithPath[Options](Root, path)
}
