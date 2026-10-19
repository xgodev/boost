package datadog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
)

func TestConfigAdd_RegistersTheKeysUnderThePath(t *testing.T) {
	path := "boost.factory.aws.clients.orders.plugins.datadog"
	ConfigAdd(path)
	config.Set(koanf.New())
	config.Load()

	assert.True(t, config.Exists(path+enabled), "%s%s is not registered", path, enabled)

	opts, err := NewOptionsWithPath(path)
	require.NoError(t, err)
	assert.True(t, opts.Enabled)
}

func TestNewOptions_EnabledByDefault(t *testing.T) {
	config.Set(koanf.New())
	config.Load()

	opts, err := NewOptions()
	require.NoError(t, err)
	assert.True(t, opts.Enabled)
}
