package main

import (
	"context"
	"os"
	"os/exec"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/bootstrap/function/middleware/logger"
	"github.com/xgodev/boost/bootstrap/function/middleware/recovery"
	"github.com/xgodev/boost/builder/examples/simple/gen/builder/worker"
	"github.com/xgodev/boost/builder/examples/simple/order"
	"github.com/xgodev/boost/builder/examples/simple/ordercreated"
	"github.com/xgodev/boost/extra/middleware"
	functionfx "github.com/xgodev/boost/fx/modules/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"go.uber.org/fx"
)

// TestMain boots boost: the function and its adapter read their config.
func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func event(t *testing.T, data any) cloudevents.Event {
	t.Helper()
	in := cloudevents.NewEvent()
	in.SetID("evt-1")
	in.SetType("order.created")
	in.SetSource("test")
	require.NoError(t, in.SetData(cloudevents.ApplicationJSON, data))
	return in
}

// TestWorker builds app worker from the generated Modules() and calls the
// function handler it provides, as the adapter does for each event.
func TestWorker(t *testing.T) {
	var (
		handle      function.Handler[*cloudevents.Event]
		svc         *order.Service
		middlewares []middleware.AnyErrorMiddleware[*cloudevents.Event]
	)
	app := fx.New(append(worker.Modules(), fx.NopLogger, fx.Populate(&handle, &svc),
		fx.Invoke(fx.Annotate(
			func(m []middleware.AnyErrorMiddleware[*cloudevents.Event]) { middlewares = m },
			fx.ParamTags(`name:"`+functionfx.MiddlewaresName+`"`),
		)),
	)...)
	require.NoError(t, app.Err())

	require.Len(t, middlewares, 2, "the function runs the middlewares of @Function")
	assert.IsType(t, &recovery.Recovery[*cloudevents.Event]{}, middlewares[0], "the outermost first")
	assert.IsType(t, &logger.Logger[*cloudevents.Event]{}, middlewares[1])

	out, err := handle(context.Background(), event(t, map[string]string{"item": "book"}))
	require.NoError(t, err)
	assert.Equal(t, ordercreated.RecordedType, out.Type())

	recorded, ok := svc.Get("1")
	require.True(t, ok, "the handler got the order service of the app")
	assert.Equal(t, order.Order{ID: "1", Item: "book"}, recorded)

	_, err = handle(context.Background(), event(t, map[string]string{}))
	assert.True(t, errors.IsNotValid(err), "got %v", err)
}

const (
	// adapterEnv chooses the adapter of the @Function of app worker
	// (worker.function.adapter, generated for adapter=[cloudevents, nats]).
	adapterEnv = "WORKER_FUNCTION_ADAPTER"
	natsURLEnv = "BOOST_FACTORY_NATS_URL"
	childEnv   = "BUILDER_EXAMPLE_CHILD"
)

// TestWorkerAdapter boots app worker in a new process per value of the
// adapter key: config is read once per process. The nats adapter connects
// to a NATS server this test runs.
func TestWorkerAdapter(t *testing.T) {
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	srv := natsserver.RunServer(&opts)
	t.Cleanup(srv.Shutdown)

	for _, value := range []string{"cloudevents", "nats", "nope"} {
		t.Run(value, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestWorkerAdapterChild$", "-test.v")
			cmd.Env = append(os.Environ(), adapterEnv+"="+value, natsURLEnv+"="+srv.ClientURL(), childEnv+"=1")
			out, err := cmd.CombinedOutput()
			require.NoError(t, err, "%s", out)
			assert.Contains(t, string(out), "--- PASS: TestWorkerAdapterChild")
		})
	}
}

// TestWorkerAdapterChild runs in the process TestWorkerAdapter starts, with
// the adapter key set.
func TestWorkerAdapterChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestWorkerAdapter")
	}

	var adapters []function.CmdFunc[*cloudevents.Event]
	app := fx.New(append(worker.Modules(), fx.NopLogger,
		fx.Invoke(fx.Annotate(
			func(a []function.CmdFunc[*cloudevents.Event]) { adapters = a },
			fx.ParamTags(`group:"`+functionfx.BSFunctionAdaptersGroupKey+`"`),
		)),
	)...)

	// the command each adapter runs the function with
	commands := map[string]string{"cloudevents": "cloudevents_http", "nats": "nats"}
	switch value := os.Getenv(adapterEnv); value {
	case "cloudevents", "nats":
		require.NoError(t, app.Err())
		require.Len(t, adapters, 1, "config wires the chosen adapter and only that one")
		assert.Equal(t, commands[value], adapters[0](nil).Use)
	default:
		require.Error(t, app.Err(), "an unknown adapter fails the boot")
		assert.True(t, errors.IsNotValid(app.Err()), "got %v", app.Err())
		assert.Contains(t, app.Err().Error(), "worker.function.adapter")
		assert.Contains(t, app.Err().Error(), `"`+value+`"`)
		assert.Contains(t, app.Err().Error(), "cloudevents, nats")
	}
}
