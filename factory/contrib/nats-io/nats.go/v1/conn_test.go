package nats

import (
	"context"
	"testing"

	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
)

func TestNewConnWithConfigPathRunsThePlugins(t *testing.T) {
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	s := natsserver.RunServer(&opts)
	defer s.Shutdown()

	const path = "boost.factory.nats.clients.test"
	ConfigAdd(path)
	t.Setenv("BOOST_FACTORY_NATS_CLIENTS_TEST_URL", s.ClientURL())
	boost.Start()

	var got *nats.Conn
	var configured bool
	plugin := func(context.Context) (func(context.Context, *nats.Conn) error, func(context.Context, []nats.Option) ([]nats.Option, error)) {
		return func(_ context.Context, c *nats.Conn) error {
				got = c
				return nil
			}, func(_ context.Context, o []nats.Option) ([]nats.Option, error) {
				configured = true
				return o, nil
			}
	}

	conn, err := NewConnWithConfigPath(context.Background(), path, plugin)
	require.NoError(t, err)
	defer conn.Close()

	assert.True(t, configured, "the plugin's options must be applied before connecting")
	assert.Same(t, conn, got, "the plugin must receive the connection")
	assert.Equal(t, s.ClientURL(), conn.ConnectedUrl(), "the url must come from the config path")
}
