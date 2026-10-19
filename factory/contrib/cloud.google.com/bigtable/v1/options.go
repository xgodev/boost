package bigtable

import "github.com/xgodev/boost/wrapper/config"

// Options holds the Bigtable client configuration.
type Options struct {
	Project     string `config:"project"`
	Instance    string `config:"instance"`
	Credentials struct {
		File string `config:"file"`
	} `config:"credentials"`
}

// NewOptions loads Options from boost.factory.bigtable.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath loads Options from the root defaults, overridden by the
// keys under path.
func NewOptionsWithPath(path string) (*Options, error) {
	return config.NewOptionsWithPath[Options](root, path)
}
