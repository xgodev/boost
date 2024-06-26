package mongo

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	"sync"

	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Manager represents a vault manager for mongodb client.
type Manager struct {
	conn    *mongo.Conn
	options *vault.ManagerOptions

	mux       sync.RWMutex
	observers map[Observer]struct{}
}

// NewManager returns a new vault manager with default options.
func NewManager(conn *mongo.Conn) *Manager {
	o, err := vault.NewManagerOptionsWithPath(root)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewManagerWithOptions(conn, o)
}

// NewManagerWithConfigPath returns a new vault manager with options from config path.
func NewManagerWithConfigPath(conn *mongo.Conn, path string) (*Manager, error) {
	o, err := vault.NewManagerOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewManagerWithOptions(conn, o), nil
}

// NewManagerWithOptions returns a new vault manager with options.
func NewManagerWithOptions(conn *mongo.Conn, options *vault.ManagerOptions) *Manager {
	return &Manager{options: options, conn: conn, observers: make(map[Observer]struct{})}
}

// Options returns vault manager options.
func (m *Manager) Options() *vault.ManagerOptions {
	return m.options
}

// Close closes mongodb client.
func (m *Manager) Close(ctx context.Context) error {
	return m.conn.Client.Disconnect(ctx)
}

// Configure configures mongodb client.
func (m *Manager) Configure(ctx context.Context, data map[string]interface{}) error {
	var username, password string
	var ok bool

	if username, ok = data["username"].(string); !ok {
		return errors.Internalf("username not found in data map")
	}

	if password, ok = data["password"].(string); !ok {
		return errors.Internalf("password not found in data map")
	}

	m.conn.Options.Auth.Username = username
	m.conn.Options.Auth.Password = password

	conn, err := mongo.NewConnWithOptions(ctx, m.conn.Options, m.conn.Plugins...)
	if err != nil {
		return err
	}

	m.Notify(conn)

	return nil
}

func (m *Manager) Register(observer Observer) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.observers[observer] = struct{}{}
}

func (m *Manager) Unregister(observer Observer) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.observers, observer)
}

func (m *Manager) Notify(conn *mongo.Conn) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	if len(m.observers) == 0 {
		log.Warn("no observers registered to receive mongo/vault notifications")
	}

	for o := range m.observers {
		o.OnNotify(conn)
	}
}
