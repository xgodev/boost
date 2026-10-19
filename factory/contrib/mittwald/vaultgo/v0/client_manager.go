package vault

import (
	"context"
	"fmt"
	"sync"
)

// ClientManager is a Manager that builds a client of type C from the
// credentials vault hands it, keeps the latest one and passes each new one
// to the functions registered with OnChange.
type ClientManager[C any] struct {
	options *ManagerOptions
	build   func(context.Context, map[string]interface{}) (C, error)
	close   func(context.Context, C) error

	mu        sync.RWMutex
	client    C
	built     bool
	observers []func(C)
}

// NewClientManager returns a manager that builds its client with build and
// closes it with close; a nil close leaves the client open.
func NewClientManager[C any](options *ManagerOptions, build func(context.Context, map[string]interface{}) (C, error), close func(context.Context, C) error) *ClientManager[C] {
	return &ClientManager[C]{options: options, build: build, close: close}
}

// Options returns the manager options.
func (m *ClientManager[C]) Options() *ManagerOptions {
	return m.options
}

// Client returns the latest client, and false before the first Configure.
func (m *ClientManager[C]) Client() (C, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client, m.built
}

// OnChange registers fn to receive every client built from now on.
func (m *ClientManager[C]) OnChange(fn func(C)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.observers = append(m.observers, fn)
}

// Configure builds a client from data, keeps it and passes it to the
// OnChange functions. A failed build keeps the previous client.
func (m *ClientManager[C]) Configure(ctx context.Context, data map[string]interface{}) error {
	client, err := m.build(ctx, data)
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.client, m.built = client, true
	observers := append([]func(C){}, m.observers...)
	m.mu.Unlock()

	for _, fn := range observers {
		fn(client)
	}
	return nil
}

// Close closes the latest client, if any.
func (m *ClientManager[C]) Close(ctx context.Context) error {
	client, built := m.Client()
	if !built || m.close == nil {
		return nil
	}
	return m.close(ctx, client)
}

// String returns data[key] as a string, or an error when it is missing or
// not a string.
func String(data map[string]interface{}, key string) (string, error) {
	v, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("vault: %s not found in data map", key)
	}
	return v, nil
}
