package elasticsearch

import (
	"context"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8"
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
)

const root = vault.ManagersRoot + ".elasticsearch"

func init() {
	vault.ManagerConfigAdd(root)
}

// Manager rebuilds an elasticsearch client with the username and password
// vault hands it.
type Manager = vault.ClientManager[*es.Client]

// NewManager returns an elasticsearch manager with default manager options.
func NewManager(o *elasticsearch.Options, plugins ...elasticsearch.Plugin) (*Manager, error) {
	return NewManagerWithConfigPath(root, o, plugins...)
}

// NewManagerWithConfigPath returns an elasticsearch manager with manager options from config path.
func NewManagerWithConfigPath(path string, o *elasticsearch.Options, plugins ...elasticsearch.Plugin) (*Manager, error) {
	mo, err := vault.NewManagerOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewManagerWithOptions(mo, o, plugins...), nil
}

// NewManagerWithOptions returns an elasticsearch manager. Each Configure
// builds a client from a copy of o with the vault credentials; o is not
// changed.
func NewManagerWithOptions(mo *vault.ManagerOptions, o *elasticsearch.Options, plugins ...elasticsearch.Plugin) *Manager {
	return vault.NewClientManager(mo,
		func(ctx context.Context, data map[string]interface{}) (*es.Client, error) {
			username, err := vault.String(data, "username")
			if err != nil {
				return nil, err
			}
			password, err := vault.String(data, "password")
			if err != nil {
				return nil, err
			}
			opts := *o
			opts.Username, opts.Password = username, password
			return elasticsearch.NewClientWithOptions(ctx, &opts, plugins...)
		},
		nil,
	)
}
