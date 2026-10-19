package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/bootstrap/function"
	fxfunction "github.com/xgodev/boost/fx/modules/bootstrap/function"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	"go.uber.org/fx"
)

// load sets the adapter keys and the factory version through env and reloads
// the config; the koanf instance is global, so every test sets what it needs.
func load(t *testing.T, enabled, version, factoryVersion string) {
	t.Helper()
	t.Setenv("BOOST_BOOTSTRAP_FUNCTION_ADAPTER_PUBSUB_ENABLED", enabled)
	t.Setenv("BOOST_BOOTSTRAP_FUNCTION_ADAPTER_PUBSUB_VERSION", version)
	t.Setenv("BOOST_FACTORY_GCP_PUBSUB_VERSION", factoryVersion)
	config.Set(koanf.New())
	config.Load()
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		enabled, version string
		want             bool
	}{
		{enabled: "true", version: "v2", want: true},
		{enabled: "true", version: "v1", want: false},
		{enabled: "true", version: "v3", want: false},
		{enabled: "false", version: "v2", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.enabled+"/"+tt.version, func(t *testing.T) {
			load(t, tt.enabled, tt.version, "v2")

			assert.Equal(t, tt.want, IsEnabled())
		})
	}
}

func TestModuleUnknownVersionFailsBoot(t *testing.T) {
	load(t, "true", "v3", "v2")

	assert.ErrorContains(t, fx.ValidateApp(Module[string]()), "boost.bootstrap.function.adapter.pubsub.version")
}

func TestModuleDisabledProvidesNothing(t *testing.T) {
	load(t, "true", "v1", "v1")

	var cmds []function.CmdFunc[string]
	app := fx.New(
		fx.NopLogger,
		Module[string](),
		fx.Invoke(fx.Annotate(
			func(c []function.CmdFunc[string]) { cmds = c },
			fx.ParamTags(`group:"`+fxfunction.BSFunctionAdaptersGroupKey+`"`),
		)),
	)
	assert.NoError(t, app.Err())
	assert.Empty(t, cmds)
}

func TestModuleGraphIsComplete(t *testing.T) {
	load(t, "true", "v2", "v2")

	assert.NoError(t, fx.ValidateApp(
		Module[string](),
		fx.Invoke(fx.Annotate(
			func([]function.CmdFunc[string]) {},
			fx.ParamTags(`group:"`+fxfunction.BSFunctionAdaptersGroupKey+`"`),
		)),
	))
}
