package mongo

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const root = vault.ManagersRoot + ".mongo"

func init() {
	vault.ManagerConfigAdd(root)
}

// Manager rebuilds a mongo connection with the username and password vault
// hands it.
type Manager = vault.ClientManager[*mongo.Conn]

// NewManager returns a mongo manager with default manager options.
func NewManager(o *mongo.Options, plugins ...mongo.Plugin) (*Manager, error) {
	return NewManagerWithConfigPath(root, o, plugins...)
}

// NewManagerWithConfigPath returns a mongo manager with manager options from config path.
func NewManagerWithConfigPath(path string, o *mongo.Options, plugins ...mongo.Plugin) (*Manager, error) {
	mo, err := vault.NewManagerOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewManagerWithOptions(mo, o, plugins...), nil
}

// NewManagerWithOptions returns a mongo manager. Each Configure connects
// with a copy of o carrying the vault credentials; o is not changed. Close
// disconnects the latest connection.
func NewManagerWithOptions(mo *vault.ManagerOptions, o *mongo.Options, plugins ...mongo.Plugin) *Manager {
	return vault.NewClientManager(mo,
		func(ctx context.Context, data map[string]interface{}) (*mongo.Conn, error) {
			username, err := vault.String(data, "username")
			if err != nil {
				return nil, err
			}
			password, err := vault.String(data, "password")
			if err != nil {
				return nil, err
			}
			opts := *o
			auth := options.Credential{}
			if o.Auth != nil {
				auth = *o.Auth
			}
			auth.Username, auth.Password = username, password
			opts.Auth = &auth
			return mongo.NewConnWithOptions(ctx, &opts, plugins...)
		},
		func(ctx context.Context, c *mongo.Conn) error { return c.Client.Disconnect(ctx) },
	)
}
