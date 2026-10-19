package nats

import (
	"os"
	"sync"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	fxfunction "github.com/xgodev/boost/fx/modules/bootstrap/function"
	"go.uber.org/fx"
)

// TestMain boots boost: the adapter reads whether it is enabled.
func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// adapters returns the option that collects the function adapters.
func adapters(got *[]function.CmdFunc[string]) fx.Option {
	return fx.Invoke(fx.Annotate(
		func(a []function.CmdFunc[string]) { *got = a },
		fx.ParamTags(`group:"`+fxfunction.BSFunctionAdaptersGroupKey+`"`),
	))
}

// reset lets a test call the modules again: they are wired once per process.
func reset(t *testing.T) {
	once = sync.Once{}
	t.Cleanup(func() { once = sync.Once{} })
}

func TestModuleWithoutClientTakesTheAppConn(t *testing.T) {
	reset(t)

	provided := false
	var got []function.CmdFunc[string]
	app := fx.New(
		fx.NopLogger,
		fx.Provide(func() *nats.Conn { provided = true; return &nats.Conn{} }),
		ModuleWithoutClient[string](),
		adapters(&got),
	)
	require.NoError(t, app.Err(), "the app provides the *nats.Conn once")

	require.Len(t, got, 1)
	assert.Equal(t, "nats", got[0](nil).Use)
	assert.True(t, provided, "the adapter reads with the app's *nats.Conn")
}

func TestModuleBringsItsConn(t *testing.T) {
	reset(t)

	var got []function.CmdFunc[string]
	assert.NoError(t, fx.ValidateApp(Module[string](), adapters(&got)), "the nats factory module provides the *nats.Conn")
}

// TestModuleWithNamedClientTakesTheNamedConn wires the adapter with the
// *nats.Conn named prices, the app's only one (#208).
func TestModuleWithNamedClientTakesTheNamedConn(t *testing.T) {
	reset(t)

	provided := false
	var got []function.CmdFunc[string]
	app := fx.New(
		fx.NopLogger,
		fx.Provide(fx.Annotate(func() *nats.Conn { provided = true; return &nats.Conn{} }, fx.ResultTags(`name:"prices"`))),
		ModuleWithNamedClient[string]("prices"),
		adapters(&got),
	)
	require.NoError(t, app.Err(), "the adapter reads with the *nats.Conn named prices")

	require.Len(t, got, 1)
	assert.Equal(t, "nats", got[0](nil).Use)
	assert.True(t, provided)
}
