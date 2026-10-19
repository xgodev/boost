package vault

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// recorder is a Manager that records what Configure receives.
type recorder struct {
	options *ManagerOptions
	mu      sync.Mutex
	got     []map[string]interface{}
	err     error
}

func (r *recorder) Options() *ManagerOptions    { return r.options }
func (r *recorder) Close(context.Context) error { return nil }
func (r *recorder) Configure(_ context.Context, d map[string]interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.got = append(r.got, d)
	return r.err
}

// fakeVault serves one secret at /v1/secret/orders and counts the reads.
func fakeVault(t *testing.T) (string, *int) {
	t.Helper()
	reads := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/secret/orders" || r.Header.Get("X-Vault-Token") != "root" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		reads++
		_, _ = w.Write([]byte(`{"lease_id":"l1","lease_duration":60,"data":{"user":"orders","pass":"s3cret"}}`))
	}))
	t.Cleanup(srv.Close)
	return srv.URL, &reads
}

func pool(t *testing.T, addr string, managers ...Manager) *ManagerPool {
	t.Helper()
	client, err := NewClientWithOptions(context.Background(), &Options{Addr: addr, Type: "TOKEN", Token: "root"})
	require.NoError(t, err)
	return NewManagerPool(client, managers...)
}

func TestConfigureMapsTheSecretKeys(t *testing.T) {
	addr, reads := fakeVault(t)
	r := &recorder{options: &ManagerOptions{
		SecretPath: "v1/secret/orders",
		Keys:       map[string]string{"user": "username", "pass": "password", "missing": "x"},
	}}

	err := pool(t, addr).Configure(context.Background(), r)

	require.NoError(t, err)
	assert.Equal(t, 1, *reads)
	assert.Equal(t, []map[string]interface{}{{"username": "orders", "password": "s3cret"}}, r.got)
}

func TestConfigureReturnsTheManagerError(t *testing.T) {
	addr, _ := fakeVault(t)
	boom := errors.New("boom")
	r := &recorder{options: &ManagerOptions{SecretPath: "v1/secret/orders"}, err: boom}

	err := pool(t, addr).Configure(context.Background(), r)

	assert.ErrorIs(t, err, boom)
}

func TestConfigureFailsOnAnUnreadableSecret(t *testing.T) {
	addr, _ := fakeVault(t)
	r := &recorder{options: &ManagerOptions{SecretPath: "v1/secret/other"}}

	err := pool(t, addr).Configure(context.Background(), r)

	assert.Error(t, err)
	assert.Empty(t, r.got)
}

func TestManageAllConfiguresEveryManager(t *testing.T) {
	addr, reads := fakeVault(t)
	a := &recorder{options: &ManagerOptions{SecretPath: "v1/secret/orders", Keys: map[string]string{"user": "u"}}}
	b := &recorder{options: &ManagerOptions{SecretPath: "v1/secret/orders", Keys: map[string]string{"pass": "p"}}}

	pool(t, addr, a, b).ManageAll(context.Background())

	assert.Equal(t, 2, *reads)
	assert.Equal(t, []map[string]interface{}{{"u": "orders"}}, a.got)
	assert.Equal(t, []map[string]interface{}{{"p": "s3cret"}}, b.got)
}

func TestNewClientWithOptionsRejectsAnUnknownType(t *testing.T) {
	_, err := NewClientWithOptions(context.Background(), &Options{Addr: "http://127.0.0.1:1", Type: "JWT"})

	assert.ErrorContains(t, err, "JWT")
}

func TestClientManagerKeepsTheLatestClientAndNotifies(t *testing.T) {
	builds := 0
	closed := ""
	m := NewClientManager(&ManagerOptions{},
		func(_ context.Context, d map[string]interface{}) (string, error) {
			builds++
			return String(d, "password")
		},
		func(_ context.Context, c string) error { closed = c; return nil },
	)
	var seen []string
	m.OnChange(func(c string) { seen = append(seen, c) })

	_, ok := m.Client()
	assert.False(t, ok)
	require.NoError(t, m.Close(context.Background()), "Close before Configure is a no-op")
	assert.Empty(t, closed)

	require.NoError(t, m.Configure(context.Background(), map[string]interface{}{"password": "a"}))
	require.NoError(t, m.Configure(context.Background(), map[string]interface{}{"password": "b"}))
	assert.Error(t, m.Configure(context.Background(), map[string]interface{}{}))

	c, ok := m.Client()
	assert.True(t, ok)
	assert.Equal(t, "b", c, "a failed build keeps the previous client")
	assert.Equal(t, []string{"a", "b"}, seen)
	assert.Equal(t, 3, builds)

	require.NoError(t, m.Close(context.Background()))
	assert.Equal(t, "b", closed)
}

func TestStringRejectsMissingAndNonStringValues(t *testing.T) {
	_, err := String(map[string]interface{}{"password": 1}, "password")
	assert.Error(t, err)
	_, err = String(map[string]interface{}{}, "password")
	assert.Error(t, err)
	v, err := String(map[string]interface{}{"password": "p"}, "password")
	require.NoError(t, err)
	assert.Equal(t, "p", v)
}

func TestNewManagerOptionsWithPathDefaultsAndEnv(t *testing.T) {
	t.Cleanup(config.Load)
	t.Setenv("APP_VAULT_ORDERS_SECRET__PATH", "v1/secret/orders")
	t.Setenv("APP_VAULT_ORDERS_WATCHER_ENABLED", "false")

	o, err := NewManagerOptionsWithPath("app.vault.orders")

	require.NoError(t, err)
	assert.Equal(t, "v1/secret/orders", o.SecretPath)
	assert.False(t, o.Watcher.Enabled)
	assert.Equal(t, 120, o.Watcher.Increment)
	assert.Equal(t, map[string]string{"username": "username", "password": "password"}, o.Keys)
}

func TestNewOptionsWithPathReadsTheEnv(t *testing.T) {
	t.Cleanup(config.Load)
	t.Setenv("APP_VAULT_ADDR", "https://vault:8200")
	t.Setenv("APP_VAULT_TYPE", "K8S")

	o, err := NewOptionsWithPath("app.vault")

	require.NoError(t, err)
	assert.Equal(t, "https://vault:8200", o.Addr)
	assert.Equal(t, "K8S", o.Type)
	assert.Equal(t, "default", o.K8s.Role)
}
