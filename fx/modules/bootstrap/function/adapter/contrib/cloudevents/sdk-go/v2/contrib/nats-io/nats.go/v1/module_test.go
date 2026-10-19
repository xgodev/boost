package nats

import (
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	fxfunction "github.com/xgodev/boost/fx/modules/bootstrap/function"
	"go.uber.org/fx"
)

const testPort = 8371

// TestMain runs a NATS server and points the boost NATS factory at it before
// boost reads the config.
func TestMain(m *testing.M) {
	opts := natsserver.DefaultTestOptions
	opts.Port = testPort
	srv := natsserver.RunServer(&opts)

	_ = os.Setenv("BOOST_FACTORY_NATS_URL", "nats://127.0.0.1:"+strconv.Itoa(testPort))
	boost.Start()

	code := m.Run()
	srv.Shutdown()
	os.Exit(code)
}

func TestModuleAddsTheCloudEventsNATSAdapter(t *testing.T) {
	var adapters []function.CmdFunc[*event.Event]
	app := fx.New(
		fx.NopLogger,
		Module[*event.Event](),
		fx.Invoke(fx.Annotate(
			func(a []function.CmdFunc[*event.Event]) { adapters = a },
			fx.ParamTags(`group:"`+fxfunction.BSFunctionAdaptersGroupKey+`"`),
		)),
	)
	require.NoError(t, app.Err())

	require.Len(t, adapters, 1)
	assert.Equal(t, "cloudevents_nats", adapters[0](nil).Use)
}

func TestModuleWithoutClientTakesTheAppConn(t *testing.T) {
	once = sync.Once{}
	t.Cleanup(func() { once = sync.Once{} })

	provided := false
	var adapters []function.CmdFunc[*event.Event]
	app := fx.New(
		fx.NopLogger,
		fx.Provide(func() *nats.Conn { provided = true; return &nats.Conn{} }),
		ModuleWithoutClient[*event.Event](),
		fx.Invoke(fx.Annotate(
			func(a []function.CmdFunc[*event.Event]) { adapters = a },
			fx.ParamTags(`group:"`+fxfunction.BSFunctionAdaptersGroupKey+`"`),
		)),
	)
	require.NoError(t, app.Err(), "the app provides the *nats.Conn once")

	require.Len(t, adapters, 1)
	assert.Equal(t, "cloudevents_nats", adapters[0](nil).Use)
	assert.True(t, provided, "the adapter reads with the app's *nats.Conn")
}
