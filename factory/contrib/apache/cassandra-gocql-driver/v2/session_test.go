package gocql

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/apache/cassandra-gocql-driver/v2"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/internal/cassandratest"
	"github.com/xgodev/boost/internal/configtest"
)

// Each plugin type must stay an alias of factory.Plugin: plugins written
// against either name keep compiling and stay interchangeable.
func TestPluginTypesAreFactoryPlugin(t *testing.T) {
	var _ factory.Plugin[*gocql.Session] = Plugin(func(context.Context, *gocql.Session) error { return nil })
}

func TestNewOptionsHasTheV1Defaults(t *testing.T) {
	boost.Start()

	o, err := NewOptions()
	require.NoError(t, err)
	require.Equal(t, &Options{
		Hosts:                    []string{"127.0.0.1"},
		Port:                     9042,
		CQLVersion:               "3.0.0",
		Timeout:                  600 * time.Millisecond,
		ConnectTimeout:           600 * time.Millisecond,
		NumConns:                 2,
		Consistency:              "QUORUM",
		MaxPreparedStmts:         1000,
		MaxRoutingKeyInfo:        1000,
		PageSize:                 5000,
		DefaultTimestamp:         true,
		ReconnectInterval:        10 * time.Millisecond,
		MaxWaitSchemaAgreement:   60 * time.Second,
		DisableInitialHostLookup: true,
		WriteCoalesceWaitTime:    200 * time.Microsecond,
	}, o)
}

func TestNewOptionsWithPathOverridesTheRoot(t *testing.T) {
	t.Setenv("APP_CASSANDRA_OVERRIDE_PORT", "19042")
	t.Setenv("APP_CASSANDRA_OVERRIDE_KEYSPACE", "orders")
	boost.Start()

	o, err := NewOptionsWithPath("app.cassandra.override")
	require.NoError(t, err)
	require.Equal(t, 19042, o.Port)
	require.Equal(t, "orders", o.Keyspace)
	require.Equal(t, []string{"127.0.0.1"}, o.Hosts, "keys the path leaves out keep the root value")
}

func TestNewClusterMapsEveryOption(t *testing.T) {
	o := &Options{
		Hosts:                    []string{"10.0.0.1", "10.0.0.2"},
		Port:                     9142,
		DC:                       "dc1",
		Username:                 "boost",
		Password:                 "secret",
		CQLVersion:               "3.4.0",
		ProtoVersion:             4,
		Timeout:                  2 * time.Second,
		ConnectTimeout:           3 * time.Second,
		Keyspace:                 "orders",
		NumConns:                 4,
		Consistency:              "local_one",
		SocketKeepalive:          30 * time.Second,
		MaxPreparedStmts:         10,
		MaxRoutingKeyInfo:        20,
		PageSize:                 100,
		DefaultTimestamp:         true,
		ReconnectInterval:        5 * time.Second,
		MaxWaitSchemaAgreement:   7 * time.Second,
		DisableInitialHostLookup: true,
		WriteCoalesceWaitTime:    time.Millisecond,
	}

	c, err := newCluster(o)
	require.NoError(t, err)
	require.Equal(t, []string{"10.0.0.1", "10.0.0.2"}, c.Hosts)
	require.Equal(t, 9142, c.Port)
	require.Equal(t, "3.4.0", c.CQLVersion)
	require.Equal(t, 4, c.ProtoVersion)
	require.Equal(t, 2*time.Second, c.Timeout)
	require.Equal(t, 3*time.Second, c.ConnectTimeout)
	require.Equal(t, "orders", c.Keyspace)
	require.Equal(t, 4, c.NumConns)
	require.Equal(t, gocql.LocalOne, c.Consistency)
	require.Equal(t, 30*time.Second, c.SocketKeepalive)
	require.Equal(t, 10, c.MaxPreparedStmts)
	require.Equal(t, 20, c.MaxRoutingKeyInfo)
	require.Equal(t, 100, c.PageSize)
	require.True(t, c.DefaultTimestamp)
	require.Equal(t, 5*time.Second, c.ReconnectInterval)
	require.Equal(t, 7*time.Second, c.MaxWaitSchemaAgreement)
	require.True(t, c.DisableInitialHostLookup)
	require.Equal(t, time.Millisecond, c.WriteCoalesceWaitTime)
	require.Equal(t, gocql.PasswordAuthenticator{Username: "boost", Password: "secret"}, c.Authenticator)
	require.NotNil(t, c.PoolConfig.HostSelectionPolicy, "a DC routes to its own nodes first")
	require.NotNil(t, c.HostFilter, "a DC filters out the nodes of other DCs")
}

func TestNewClusterKeepsTheDriverDefaultsForZeroOptions(t *testing.T) {
	want := gocql.NewCluster("10.0.0.1")

	c, err := newCluster(&Options{Hosts: []string{"10.0.0.1"}})
	require.NoError(t, err)
	require.Equal(t, want.Port, c.Port)
	require.Equal(t, want.CQLVersion, c.CQLVersion)
	require.Equal(t, want.ProtoVersion, c.ProtoVersion)
	require.Equal(t, want.Timeout, c.Timeout)
	require.Equal(t, want.ConnectTimeout, c.ConnectTimeout)
	require.Equal(t, want.Keyspace, c.Keyspace)
	require.Equal(t, want.NumConns, c.NumConns)
	require.Equal(t, want.Consistency, c.Consistency)
	require.Equal(t, want.SocketKeepalive, c.SocketKeepalive)
	require.Equal(t, want.MaxPreparedStmts, c.MaxPreparedStmts)
	require.Equal(t, want.MaxRoutingKeyInfo, c.MaxRoutingKeyInfo)
	require.Equal(t, want.PageSize, c.PageSize)
	require.Equal(t, want.ReconnectInterval, c.ReconnectInterval)
	require.Equal(t, want.MaxWaitSchemaAgreement, c.MaxWaitSchemaAgreement)
	require.Equal(t, want.WriteCoalesceWaitTime, c.WriteCoalesceWaitTime)
	require.False(t, c.DefaultTimestamp, "the option, not the driver, decides the client timestamp")
	require.False(t, c.DisableInitialHostLookup, "the option, not the driver, decides the host lookup")
	require.Nil(t, c.Authenticator, "no credentials, no authenticator")
	require.Nil(t, c.PoolConfig.HostSelectionPolicy, "no DC, the driver's policy")
	require.Nil(t, c.HostFilter, "no DC, no filter")
}

func TestNewClusterAuthenticatesWithAPasswordOnly(t *testing.T) {
	c, err := newCluster(&Options{Password: "secret"})
	require.NoError(t, err)
	require.Equal(t, gocql.PasswordAuthenticator{Password: "secret"}, c.Authenticator)
}

func TestNewClusterRejectsAnUnknownConsistency(t *testing.T) {
	_, err := newCluster(&Options{Consistency: "MOST"})
	require.ErrorContains(t, err, "MOST")
}

func TestNewSessionWithOptionsRejectsAnUnknownConsistency(t *testing.T) {
	s, err := NewSessionWithOptions(context.Background(), &Options{Hosts: []string{"127.0.0.1"}, Consistency: "MOST"})
	require.ErrorContains(t, err, "MOST")
	require.Nil(t, s)
}

func TestNewSessionWithOptionsFailsWhenCassandraIsUnreachable(t *testing.T) {
	port := cassandratest.ClosedPort(t)

	s, err := NewSessionWithOptions(context.Background(), &Options{
		Hosts:          []string{"127.0.0.1"},
		Port:           port,
		ConnectTimeout: 200 * time.Millisecond,
	})
	require.ErrorContains(t, err, fmt.Sprintf("127.0.0.1:%d", port))
	require.Nil(t, s)
}

func TestNewSessionFailsOnABadConfigValue(t *testing.T) {
	configtest.Set(t, "BOOST_FACTORY_GOCQL_PORT", "ninety", "9042")

	s, err := NewSession(context.Background())
	require.Error(t, err)
	require.Nil(t, s)
}

func TestNewSessionWithConfigPathFailsOnABadConfigValue(t *testing.T) {
	t.Setenv("APP_CASSANDRA_BAD_PORT", "ninety")
	boost.Start()

	s, err := NewSessionWithConfigPath(context.Background(), "app.cassandra.bad")
	require.Error(t, err)
	require.Nil(t, s)
}

func TestNewSessionWithConfigPathDialsThePathAddress(t *testing.T) {
	port := cassandratest.ClosedPort(t)
	t.Setenv("APP_CASSANDRA_DIAL_PORT", strconv.Itoa(port))
	t.Setenv("APP_CASSANDRA_DIAL_CONNECTTIMEOUT", "200ms")
	boost.Start()

	s, err := NewSessionWithConfigPath(context.Background(), "app.cassandra.dial")
	require.ErrorContains(t, err, fmt.Sprintf("127.0.0.1:%d", port))
	require.Nil(t, s)
}

func liveOptions(host string, port int) *Options {
	return &Options{
		Hosts:                    []string{host},
		Port:                     port,
		Timeout:                  10 * time.Second,
		ConnectTimeout:           10 * time.Second,
		Consistency:              "ONE",
		DisableInitialHostLookup: true,
	}
}

func releaseVersion(t *testing.T, s *gocql.Session) string {
	t.Helper()
	var v string
	require.NoError(t, s.Query("SELECT release_version FROM system.local").Scan(&v))
	return v
}

func TestNewSessionWithOptionsConnectsAndAppliesThePlugins(t *testing.T) {
	host, port := cassandratest.Addr(t)

	var got []*gocql.Session
	plugin := func(_ context.Context, s *gocql.Session) error {
		got = append(got, s)
		return nil
	}

	s, err := NewSessionWithOptions(context.Background(), liveOptions(host, port), plugin, plugin)
	require.NoError(t, err)
	t.Cleanup(s.Close)

	require.Equal(t, []*gocql.Session{s, s}, got, "every plugin sees the session, in order")
	require.NotEmpty(t, releaseVersion(t, s))
}

func TestNewSessionWithOptionsClosesTheSessionWhenAPluginFails(t *testing.T) {
	host, port := cassandratest.Addr(t)

	boom := errors.New("boom")
	var got *gocql.Session
	s, err := NewSessionWithOptions(context.Background(), liveOptions(host, port), func(_ context.Context, s *gocql.Session) error {
		got = s
		return boom
	})
	require.ErrorIs(t, err, boom)
	require.Nil(t, s)
	require.NotNil(t, got)
	require.True(t, got.Closed(), "a failed plugin leaves no open session behind")
}

func TestNewSessionConnectsWithTheConfig(t *testing.T) {
	host, port := cassandratest.Addr(t)
	t.Setenv("BOOST_FACTORY_GOCQL_HOSTS", host)
	t.Setenv("BOOST_FACTORY_GOCQL_PORT", strconv.Itoa(port))
	t.Setenv("BOOST_FACTORY_GOCQL_CONSISTENCY", "ONE")
	t.Setenv("BOOST_FACTORY_GOCQL_TIMEOUT", "10s")
	t.Setenv("BOOST_FACTORY_GOCQL_CONNECTTIMEOUT", "10s")
	boost.Start()

	s, err := NewSession(context.Background())
	require.NoError(t, err)
	t.Cleanup(s.Close)

	require.NotEmpty(t, releaseVersion(t, s))
}
