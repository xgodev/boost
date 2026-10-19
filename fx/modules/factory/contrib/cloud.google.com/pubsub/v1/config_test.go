package pubsub

import (
	"testing"

	gpubsub "cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	"go.uber.org/fx"
)

// load sets the version through env and reloads the config; the koanf
// instance is global, so every test sets the value it needs.
func load(t *testing.T, version string) {
	t.Helper()
	t.Setenv("BOOST_FACTORY_GCP_PUBSUB_VERSION", version)
	config.Set(koanf.New())
	config.Load()
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{version: "v1", want: true},
		{version: "v2", want: false},
		{version: "v3", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			load(t, tt.version)

			assert.Equal(t, tt.want, IsEnabled())
		})
	}
}

func TestModuleUnknownVersionFailsBoot(t *testing.T) {
	load(t, "v3")

	assert.ErrorContains(t, fx.ValidateApp(Module()), "boost.factory.gcp.pubsub.version")
}

func TestModuleDisabledProvidesNothing(t *testing.T) {
	load(t, "v2")

	var client *gpubsub.Client
	assert.ErrorContains(t, fx.ValidateApp(Module(), fx.Populate(&client)), "missing type")
}

func TestModuleGraphIsComplete(t *testing.T) {
	load(t, "v1")

	assert.NoError(t, fx.ValidateApp(Module()))
}
