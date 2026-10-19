package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	h "github.com/xgodev/boost/extra/health"
	factory "github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v9"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
)

// connect returns a client to a fake node that answers the creation ping
// with 200 and every later request with status.
func connect(t *testing.T, status int) *elasticsearch.Client {
	t.Helper()
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)
	client, err := factory.NewClientWithOptions(context.Background(), &factory.Options{Addresses: []string{srv.URL}, DisableRetry: true})
	require.NoError(t, err)
	return client
}

// loadEnv sets key and loads config; cleanup reloads it with restore, since
// boost's koanf keeps keys across Load calls.
func loadEnv(t *testing.T, key, value, restore string) {
	t.Helper()
	t.Setenv(key, value)
	config.Set(koanf.New())
	config.Load()
	t.Cleanup(func() {
		_ = os.Setenv(key, restore) // intentional: t.Setenv already proved the key is settable
		config.Load()
	})
}

// result finds the check named name in the health package.
func result(t *testing.T, name string) *h.CheckResult {
	t.Helper()
	for _, r := range h.CheckAll(context.Background()) {
		if r.HealthCheck.Name == name {
			return r
		}
	}
	return nil
}

func TestRegisterAddsACheckerThatPingsTheNode(t *testing.T) {
	loadEnv(t, "BOOST_FACTORY_ELASTICSEARCH_PLUGINS_HEALTH_NAME", "es-v9-up", "elasticsearch")

	require.NoError(t, Register(context.Background(), connect(t, http.StatusOK)))

	r := result(t, "es-v9-up")
	require.NotNil(t, r, "the checker is registered under the configured name")
	assert.True(t, r.IsOk())
	assert.True(t, r.HealthCheck.IsRequired())
}

func TestRegisterFailsOnUnreadableConfig(t *testing.T) {
	loadEnv(t, "BOOST_FACTORY_ELASTICSEARCH_PLUGINS_HEALTH_REQUIRED", "maybe", "true")

	assert.Error(t, Register(context.Background(), nil))
}

func TestCheckerReportsARejectedPing(t *testing.T) {
	err := NewChecker(connect(t, http.StatusServiceUnavailable)).Check(context.Background())

	require.Error(t, err)
	assert.True(t, errors.IsServiceUnavailable(err), "got %v", err)
	assert.Contains(t, err.Error(), "503")
}

func TestCheckerReportsAnUnreachableNode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
	}))
	client, err := factory.NewClientWithOptions(context.Background(), &factory.Options{Addresses: []string{srv.URL}, DisableRetry: true})
	require.NoError(t, err)
	srv.Close()

	assert.Error(t, NewChecker(client).Check(context.Background()))
}

func TestNewHealthReadsTheDefaults(t *testing.T) {
	config.Set(koanf.New())
	config.Load()

	got, err := NewHealth()

	require.NoError(t, err)
	assert.Equal(t, &Options{Name: "elasticsearch", Description: "default connection", Required: true, Enabled: true}, got.options)
}

func TestNewHealthFailsOnUnreadableConfig(t *testing.T) {
	loadEnv(t, "BOOST_FACTORY_ELASTICSEARCH_PLUGINS_HEALTH_ENABLED", "maybe", "true")

	got, err := NewHealth()

	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestNewHealthWithConfigPathReadsTheGivenPath(t *testing.T) {
	loadEnv(t, "APP_SEARCH_HEALTH_NAME", "search", "")

	got, err := NewHealthWithConfigPath("app.search.health")

	require.NoError(t, err)
	assert.Equal(t, "search", got.options.Name)
	assert.Equal(t, "default connection", got.options.Description, "what the path leaves out comes from the plugin root")
}

func TestNewHealthWithConfigPathFailsOnUnreadablePath(t *testing.T) {
	loadEnv(t, "APP_BROKEN_HEALTH_REQUIRED", "maybe", "true")

	got, err := NewHealthWithConfigPath("app.broken.health")

	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestADisabledHealthIsNotChecked(t *testing.T) {
	health := NewHealthWithOptions(&Options{Name: "es-v9-off", Enabled: false})

	require.NoError(t, health.Register(context.Background(), connect(t, http.StatusOK)))

	assert.Nil(t, result(t, "es-v9-off"))
}
