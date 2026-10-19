package vault

import "github.com/xgodev/boost/wrapper/config"

// ManagerOptions vault manager options.
type ManagerOptions struct {
	// SecretPath is the full vault path, e.g. v1/database/creds/orders.
	SecretPath string `config:"secretPath"`
	Watcher    struct {
		Enabled   bool
		Increment int
	}
	// Keys maps secret data keys to the keys the manager reads.
	Keys map[string]string
}

// NewManagerOptionsWithPath unmarshals the manager options under path.
func NewManagerOptionsWithPath(path string) (opts *ManagerOptions, err error) {
	ManagerConfigAdd(path)
	config.Load()
	return config.NewOptionsWithPath[ManagerOptions](path)
}
