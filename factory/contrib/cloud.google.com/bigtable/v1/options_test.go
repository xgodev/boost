package bigtable

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
)

func TestNewOptions_Defaults(t *testing.T) {
	config.Set(koanf.New())
	config.Load()

	o, err := NewOptions()
	require.NoError(t, err)
	assert.Equal(t, &Options{}, o, "no project, instance or credentials file: ADC and the app's own keys")
}

func TestNewOptions_FromEnv(t *testing.T) {
	loadWith(t, "BOOST_FACTORY_BIGTABLE_PROJECT", "my-project", "")
	loadWith(t, "BOOST_FACTORY_BIGTABLE_INSTANCE", "my-instance", "")
	loadWith(t, "BOOST_FACTORY_BIGTABLE_CREDENTIALS_FILE", "/secrets/sa.json", "")

	o, err := NewOptions()
	require.NoError(t, err)
	assert.Equal(t, "my-project", o.Project)
	assert.Equal(t, "my-instance", o.Instance)
	assert.Equal(t, "/secrets/sa.json", o.Credentials.File)
}

func TestNewOptionsWithPath_NamedClient(t *testing.T) {
	path := root + ".clients.orders"
	ConfigAdd(path)
	loadWith(t, "BOOST_FACTORY_BIGTABLE_CLIENTS_ORDERS_PROJECT", "orders-project", "")
	loadWith(t, "BOOST_FACTORY_BIGTABLE_CLIENTS_ORDERS_INSTANCE", "orders-instance", "")

	o, err := NewOptionsWithPath(path)
	require.NoError(t, err)
	assert.Equal(t, "orders-project", o.Project)
	assert.Equal(t, "orders-instance", o.Instance)

	def, err := NewOptions()
	require.NoError(t, err)
	assert.Empty(t, def.Instance, "the named client's keys stay under its own path")
}
