package health

import (
	"context"
	"testing"
	"time"

	"github.com/apache/cassandra-gocql-driver/v2"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/extra/health"
	factory "github.com/xgodev/boost/factory/contrib/apache/cassandra-gocql-driver/v2"
	"github.com/xgodev/boost/internal/cassandratest"
	"github.com/xgodev/boost/internal/configtest"
)

func TestNewHealthReadsTheDefaults(t *testing.T) {
	boost.Start()

	h, err := NewHealth()
	require.NoError(t, err)
	require.Equal(t, &Options{Name: "cassandra", Description: "default connection", Required: true, Enabled: true}, h.options)
}

func TestNewHealthFailsOnABadConfigValue(t *testing.T) {
	configtest.Set(t, "BOOST_FACTORY_GOCQL_PLUGINS_HEALTH_REQUIRED", "maybe", "true")

	h, err := NewHealth()
	require.Error(t, err)
	require.Nil(t, h)
}

func TestNewHealthWithConfigPathOverridesTheRoot(t *testing.T) {
	t.Setenv("APP_CASSANDRA_OVERRIDE_NAME", "orders")
	t.Setenv("APP_CASSANDRA_OVERRIDE_REQUIRED", "false")
	boost.Start()

	h, err := NewHealthWithConfigPath("app.cassandra.override")
	require.NoError(t, err)
	require.Equal(t, &Options{Name: "orders", Description: "default connection", Required: false, Enabled: true}, h.options)
}

func TestNewHealthWithConfigPathFailsOnABadConfigValue(t *testing.T) {
	t.Setenv("APP_CASSANDRA_BAD_REQUIRED", "maybe")
	boost.Start()

	h, err := NewHealthWithConfigPath("app.cassandra.bad")
	require.Error(t, err)
	require.Nil(t, h)
}

func TestRegisterFailsOnABadConfigValue(t *testing.T) {
	configtest.Set(t, "BOOST_FACTORY_GOCQL_PLUGINS_HEALTH_REQUIRED", "maybe", "true")

	require.Error(t, Register(context.Background(), nil), "a broken health config fails the plugin, it is not skipped")
}

func liveSession(t *testing.T) *gocql.Session {
	t.Helper()
	host, port := cassandratest.Addr(t)
	s, err := factory.NewSessionWithOptions(context.Background(), &factory.Options{
		Hosts:                    []string{host},
		Port:                     port,
		Timeout:                  10 * time.Second,
		ConnectTimeout:           10 * time.Second,
		Consistency:              "ONE",
		DisableInitialHostLookup: true,
	})
	require.NoError(t, err)
	t.Cleanup(s.Close)
	return s
}

func result(t *testing.T, name string) *health.CheckResult {
	t.Helper()
	for _, r := range health.CheckAll(context.Background()) {
		if r.HealthCheck.Name == name {
			return r
		}
	}
	t.Fatalf("no health check named %q", name)
	return nil
}

func TestCheckerIsUpWhileTheSessionAnswers(t *testing.T) {
	s := liveSession(t)

	require.NoError(t, NewChecker(s).Check(context.Background()))
}

func TestCheckerIsDownOnceTheSessionIsClosed(t *testing.T) {
	s := liveSession(t)
	s.Close()

	require.Error(t, NewChecker(s).Check(context.Background()))
}

func TestRegisterAddsTheCheckerWithTheConfigOptions(t *testing.T) {
	s := liveSession(t)
	configtest.Set(t, "BOOST_FACTORY_GOCQL_PLUGINS_HEALTH_NAME", "cassandra-from-config", "cassandra")

	require.NoError(t, Register(context.Background(), s))

	r := result(t, "cassandra-from-config")
	require.True(t, r.IsOk(), "check failed: %v", r.Error)
	require.True(t, r.HealthCheck.IsRequired())
}

func TestHealthRegisterAddsTheCheckerWithItsOptions(t *testing.T) {
	s := liveSession(t)
	boost.Start()

	h := NewHealthWithOptions(&Options{Name: "cassandra-optional", Description: "reports", Required: false, Enabled: true})
	require.NoError(t, h.Register(context.Background(), s))

	r := result(t, "cassandra-optional")
	require.True(t, r.IsOk(), "check failed: %v", r.Error)
	require.False(t, r.HealthCheck.IsRequired())
	require.Equal(t, "reports", r.HealthCheck.Description)
}
