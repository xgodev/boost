package vault

import (
	"github.com/xgodev/boost/wrapper/config"
)

// ManagerOptions represents a vault client options.
type ManagerOptions struct {
	SecretPath string
	Watcher    struct {
		Enabled   bool
		Increment int
	}
	Keys map[string]string
}

// NewManagerOptionsWithPath unmarshals manager options based a given key path.
func NewManagerOptionsWithPath(path string) (opts *ManagerOptions, err error) {
	return config.NewOptionsWithPath[ManagerOptions](path)
}
