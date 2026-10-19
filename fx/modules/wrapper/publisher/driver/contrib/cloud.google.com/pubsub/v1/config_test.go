package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	"github.com/xgodev/boost/wrapper/publisher"
	"go.uber.org/fx"
)

// load sets the driver keys and the factory version through env and reloads
// the config; the koanf instance is global, so every test sets what it needs.
func load(t *testing.T, enabled, version, factoryVersion string) {
	t.Helper()
	t.Setenv("BOOST_WRAPPER_PUBLISHER_DRIVER_PUBSUB_ENABLED", enabled)
	t.Setenv("BOOST_WRAPPER_PUBLISHER_DRIVER_PUBSUB_VERSION", version)
	t.Setenv("BOOST_FACTORY_GCP_PUBSUB_VERSION", factoryVersion)
	config.Set(koanf.New())
	config.Load()
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		enabled, version string
		want             bool
	}{
		{enabled: "true", version: "v1", want: true},
		{enabled: "true", version: "v2", want: false},
		{enabled: "true", version: "v3", want: false},
		{enabled: "false", version: "v1", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.enabled+"/"+tt.version, func(t *testing.T) {
			load(t, tt.enabled, tt.version, "v1")

			assert.Equal(t, tt.want, IsEnabled())
		})
	}
}

func TestModuleUnknownVersionFailsBoot(t *testing.T) {
	load(t, "true", "v3", "v1")

	assert.ErrorContains(t, fx.ValidateApp(Module()), "boost.wrapper.publisher.driver.pubsub.version")
}

func TestModuleDisabledProvidesNothing(t *testing.T) {
	load(t, "true", "v2", "v2")

	var drv publisher.Driver
	assert.ErrorContains(t, fx.ValidateApp(Module(), fx.Populate(&drv)), "missing type")
}

func TestModuleGraphIsComplete(t *testing.T) {
	load(t, "true", "v1", "v1")

	var drv publisher.Driver
	assert.NoError(t, fx.ValidateApp(Module(), fx.Populate(&drv)))
}
