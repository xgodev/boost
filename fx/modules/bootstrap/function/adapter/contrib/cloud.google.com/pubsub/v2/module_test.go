package pubsub

import (
	"sync"
	"testing"

	"cloud.google.com/go/pubsub/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/bootstrap/function"
	fxfunction "github.com/xgodev/boost/fx/modules/bootstrap/function"
	"go.uber.org/fx"
)

func TestModuleWithoutClientTakesTheAppClient(t *testing.T) {
	load(t, "true", "v2", "v2")
	once = sync.Once{}
	t.Cleanup(func() { once = sync.Once{} })

	provided := false
	var cmds []function.CmdFunc[string]
	app := fx.New(
		fx.NopLogger,
		fx.Provide(func() *pubsub.Client { provided = true; return &pubsub.Client{} }),
		ModuleWithoutClient[string](),
		fx.Invoke(fx.Annotate(
			func(c []function.CmdFunc[string]) { cmds = c },
			fx.ParamTags(`group:"`+fxfunction.BSFunctionAdaptersGroupKey+`"`),
		)),
	)
	require.NoError(t, app.Err(), "the app provides the *pubsub.Client once")

	require.Len(t, cmds, 1)
	assert.True(t, provided, "the adapter reads with the app's *pubsub.Client")
}
