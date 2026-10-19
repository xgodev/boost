package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/pstest"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	natsserver "github.com/nats-io/nats-server/v2/test"
	natsgo "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	natsadapter "github.com/xgodev/boost/bootstrap/function/adapter/contrib/nats-io/nats.go/v1"
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
	// (worker.function.adapter, generated for adapter=[cloudevents, nats, pubsub]).
	adapterEnv = "WORKER_FUNCTION_ADAPTER"
	// driverEnv chooses the driver of the publisher the handler announces
	// with (boost.wrapper.publisher.driver.name).
	driverEnv   = "BOOST_WRAPPER_PUBLISHER_DRIVER_NAME"
	natsURLEnv  = "BOOST_FACTORY_NATS_URL"
	projectEnv  = "BOOST_FACTORY_GCP_PUBSUB_API__OPTIONS_PROJECT__ID"
	emulatorEnv = "PUBSUB_EMULATOR_HOST"
	childEnv    = "BUILDER_EXAMPLE_CHILD"
)

// TestWorkerAdapter boots app worker in a new process per adapter and
// publisher driver: config is read once per process. The function reads
// from the broker the handler announces to, or from another one; the nats
// and pubsub clients connect to a NATS server and a Pub/Sub fake this test
// runs.
func TestWorkerAdapter(t *testing.T) {
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	srv := natsserver.RunServer(&opts)
	t.Cleanup(srv.Shutdown)

	fake := pstest.NewServer()
	t.Cleanup(func() { _ = fake.Close() })

	for _, c := range []struct{ adapter, driver string }{
		{"cloudevents", "noop"},
		{"nats", "noop"},
		{"pubsub", "noop"},
		{"nats", "nats"},
		{"pubsub", "pubsub"},
		{"nats", "pubsub"},
		{"pubsub", "nats"},
		{"nope", "noop"},
	} {
		t.Run(c.adapter+"+"+c.driver, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestWorkerAdapterChild$", "-test.v")
			cmd.Env = append(os.Environ(), adapterEnv+"="+c.adapter, driverEnv+"="+c.driver,
				natsURLEnv+"="+srv.ClientURL(), emulatorEnv+"="+fake.Addr, projectEnv+"=example", childEnv+"=1")
			out, err := cmd.CombinedOutput()
			require.NoError(t, err, "%s", out)
			assert.Contains(t, string(out), "--- PASS: TestWorkerAdapterChild")
		})
	}
}

// mainEnv makes the process the TestMain_* tests start run the generated
// main instead of the tests.
const mainEnv = "BUILDER_EXAMPLE_MAIN"

// TestMain_BadConfigExitsWithTheReason runs the generated main of app worker
// with an unknown adapter and every log level at its default: the process
// exits non-zero and says why, so a pod that restarts shows its reason.
func TestMain_BadConfigExitsWithTheReason(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestMainChild$")
	cmd.Env = append(os.Environ(), adapterEnv+"=nope", mainEnv+"=1")
	out, err := cmd.CombinedOutput()

	var exit *exec.ExitError
	require.ErrorAs(t, err, &exit, "the main must fail: %s", out)
	assert.NotZero(t, exit.ExitCode())
	assert.Contains(t, string(out), "worker.function.adapter")
	assert.Contains(t, string(out), `unknown adapter "nope"`)
}

// TestMain_RunsTheChosenAdapter runs the generated main of app worker with
// only WORKER_FUNCTION_ADAPTER=nats set (#203): no CLI argument and no
// boost.bootstrap.function.default. The worker subscribes and records the
// order this test publishes.
func TestMain_RunsTheChosenAdapter(t *testing.T) {
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	srv := natsserver.RunServer(&opts)
	t.Cleanup(srv.Shutdown)

	conn, err := natsgo.Connect(srv.ClientURL())
	require.NoError(t, err)
	t.Cleanup(conn.Close)

	var out syncBuffer
	cmd := exec.Command(os.Args[0], "-test.run=^TestMainChild$")
	cmd.Env = append(os.Environ(), adapterEnv+"=nats", natsURLEnv+"="+srv.ClientURL(), mainEnv+"=1")
	cmd.Stdout, cmd.Stderr = &out, &out
	require.NoError(t, cmd.Start())
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	adapter, err := natsadapter.DefaultOptions()
	require.NoError(t, err)
	require.NotEmpty(t, adapter.Subjects)
	payload, err := json.Marshal(event(t, map[string]string{"item": "book"}))
	require.NoError(t, err)

	// the worker subscribes once it boots: publish until it records the order
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	deadline := time.After(10 * time.Second)
	for !strings.Contains(out.String(), "recorded order 1 of book") {
		select {
		case <-tick.C:
			require.NoError(t, conn.Publish(adapter.Subjects[0], payload))
		case err := <-exited:
			t.Fatalf("the worker exited (%v) without handling the order: %s", err, out.String())
		case <-deadline:
			t.Fatalf("the worker never handled the order: %s", out.String())
		}
	}
}

// syncBuffer is a bytes.Buffer the child process and the test share.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// TestMainChild runs the generated main, with no argument as a deployed
// binary gets, in the process the TestMain_* tests start.
func TestMainChild(t *testing.T) {
	if os.Getenv(mainEnv) != "1" {
		t.Skip("run by the TestMain_* tests")
	}
	os.Args = os.Args[:1]
	main()
}

// TestWorkerAdapterChild runs in the process TestWorkerAdapter starts, with
// the adapter and driver keys set. The adapter and the publisher share the
// one default client of their broker.
func TestWorkerAdapterChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestWorkerAdapter")
	}

	var (
		adapters []function.CmdFunc[*cloudevents.Event]
		handle   function.Handler[*cloudevents.Event]
	)
	app := fx.New(append(worker.Modules(), fx.NopLogger, fx.Populate(&handle),
		fx.Invoke(fx.Annotate(
			func(a []function.CmdFunc[*cloudevents.Event]) { adapters = a },
			fx.ParamTags(`group:"`+functionfx.BSFunctionAdaptersGroupKey+`"`),
		)),
	)...)

	// the command each adapter runs the function with
	commands := map[string]string{"cloudevents": "cloudevents_http", "nats": "nats", "pubsub": "gcp_pubsub"}
	switch value := os.Getenv(adapterEnv); value {
	case "cloudevents", "nats", "pubsub":
		require.NoError(t, app.Err())
		require.Len(t, adapters, 1, "config wires the chosen adapter and only that one")
		assert.Equal(t, commands[value], adapters[0](nil).Use)
		assert.NotNil(t, handle, "the handler got the publisher of the chosen driver")
	default:
		require.Error(t, app.Err(), "an unknown adapter fails the boot")
		assert.True(t, errors.IsNotValid(app.Err()), "got %v", app.Err())
		assert.Contains(t, app.Err().Error(), "worker.function.adapter")
		assert.Contains(t, app.Err().Error(), `"`+value+`"`)
		assert.Contains(t, app.Err().Error(), "cloudevents, nats, pubsub")
	}
}
