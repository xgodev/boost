package hystrix

import (
	"github.com/xgodev/boost/wrapper/config"
	"strings"
)

// Options struct which represents options.
type Options struct {
	Enabled bool
}

// NewOptions returns options from config path.
func NewOptions(name string) (opts *Options, err error) {
	opts = &Options{}
	path := strings.Join([]string{root, ".", name}, "")

	return config.NewOptionsWithPath[Options](path)
}
