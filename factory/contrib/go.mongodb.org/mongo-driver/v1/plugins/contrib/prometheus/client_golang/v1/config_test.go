package prometheus

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
)

// mongoKey is where the default prometheus plugin of the mongo factory
// keeps its switch: under the mongo factory's plugins root.
const mongoKey = "boost.factory.mongo.plugins.prometheus.enabled"

func TestMain(m *testing.M) {
	_ = os.Setenv("BOOST_FACTORY_MONGO_PLUGINS_PROMETHEUS_ENABLED", "false")
	boost.Start()
	os.Exit(m.Run())
}

func TestDefaultKeyIsUnderTheMongoFactory(t *testing.T) {
	assert.True(t, config.Exists(mongoKey), "the default plugin registers %s", mongoKey)
	assert.False(t, config.Exists("boost.factory.resty.plugins.prometheus.enabled"),
		"the mongo plugin registers nothing under the resty factory")

	opts, err := NewOptions()
	require.NoError(t, err)
	assert.False(t, opts.Enabled, "the mongo key turns the default plugin off")
}
