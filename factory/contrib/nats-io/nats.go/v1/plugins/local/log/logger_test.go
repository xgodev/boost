package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	ginats "github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1"
)

const testPath = "boost.factory.nats.clients.test.plugins.log"

var _ ginats.Plugin = (&Log{}).Register

func TestLogRegisterAddsTheConnectionHandlers(t *testing.T) {
	ConfigAdd(testPath)
	boost.Start()

	l, err := NewLogWithConfigPath(testPath)
	require.NoError(t, err)

	connPlugin, optsPlugin := l.Register(context.Background())
	assert.Nil(t, connPlugin)
	require.NotNil(t, optsPlugin)

	opts, err := optsPlugin(context.Background(), nil)
	require.NoError(t, err)
	assert.Len(t, opts, 3, "disconnect, reconnect and closed handlers")
}

func TestLogRegisterDoesNothingWhenDisabled(t *testing.T) {
	ConfigAdd(testPath)
	t.Setenv("BOOST_FACTORY_NATS_CLIENTS_TEST_PLUGINS_LOG_ENABLED", "false")
	boost.Start()

	l, err := NewLogWithConfigPath(testPath)
	require.NoError(t, err)

	connPlugin, optsPlugin := l.Register(context.Background())
	assert.Nil(t, connPlugin)
	assert.Nil(t, optsPlugin)
}
