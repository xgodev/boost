package gocql

import (
	"context"
	"strconv"
	"testing"

	"github.com/apache/cassandra-gocql-driver/v2"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	g "github.com/xgodev/boost/factory/contrib/apache/cassandra-gocql-driver/v2"
	"github.com/xgodev/boost/internal/cassandratest"
	"github.com/xgodev/boost/internal/configtest"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModuleGraphIsComplete(t *testing.T) {
	boost.Start()

	first, second := Module(), Module()
	require.NoError(t, fx.ValidateApp(first, second, fx.Invoke(func(*gocql.Session) {})),
		"the first Module() provides the session and a second adds nothing, so dependents may each include it")
}

func TestModuleFailsTheBootWhenCassandraIsUnreachable(t *testing.T) {
	configtest.Set(t, "BOOST_FACTORY_GOCQL_PORT", strconv.Itoa(cassandratest.ClosedPort(t)), "9042")
	configtest.Set(t, "BOOST_FACTORY_GOCQL_CONNECTTIMEOUT", "200ms", "600ms")

	app := fx.New(fx.NopLogger, fx.Provide(context.Background), module(), fx.Invoke(func(*gocql.Session) {}))
	require.Error(t, app.Err())
}

func TestModuleAppliesThePluginsAndClosesTheSessionOnStop(t *testing.T) {
	host, port := cassandratest.Addr(t)
	t.Setenv("BOOST_FACTORY_GOCQL_HOSTS", host)
	t.Setenv("BOOST_FACTORY_GOCQL_PORT", strconv.Itoa(port))
	t.Setenv("BOOST_FACTORY_GOCQL_CONSISTENCY", "ONE")
	t.Setenv("BOOST_FACTORY_GOCQL_TIMEOUT", "10s")
	t.Setenv("BOOST_FACTORY_GOCQL_CONNECTTIMEOUT", "10s")
	boost.Start()

	var applied *gocql.Session
	plugins := []g.Plugin{func(_ context.Context, s *gocql.Session) error {
		applied = s
		return nil
	}}

	var s *gocql.Session
	app := fxtest.New(t, fx.Provide(context.Background), fx.Supply(plugins), module(), fx.Populate(&s))
	app.RequireStart()
	require.Same(t, s, applied, "the plugins in the graph see the session")
	require.False(t, s.Closed())

	app.RequireStop()
	require.True(t, s.Closed(), "the session is closed on app stop")
}
