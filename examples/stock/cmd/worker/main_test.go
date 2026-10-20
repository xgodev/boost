package main

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	cehttp "github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http"
	"github.com/xgodev/boost/bootstrap/function/middleware/logger"
	"github.com/xgodev/boost/bootstrap/function/middleware/recovery"
	"github.com/xgodev/boost/examples/stock/gen/builder/worker"
	"github.com/xgodev/boost/examples/stock/gen/implementation/internal_/pkg/domain/service/stock/reserver"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/repository/stock"
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
// function handler it provides, as the adapter does for each event. The
// reserver is local: it takes from the memory stock.
func TestWorker(t *testing.T) {
	var (
		handle      function.Handler[[]*cloudevents.Event]
		repo        stock.Repository
		middlewares []middleware.AnyErrorMiddleware[[]*cloudevents.Event]
	)
	app := fx.New(append(worker.Modules(), fx.NopLogger, fx.Populate(&handle, &repo),
		fx.Invoke(fx.Annotate(
			func(m []middleware.AnyErrorMiddleware[[]*cloudevents.Event]) { middlewares = m },
			fx.ParamTags(`name:"`+functionfx.MiddlewaresName+`"`),
		)),
	)...)
	require.NoError(t, app.Err())

	require.Len(t, middlewares, 2, "the function runs the middlewares of @Function")
	assert.IsType(t, &recovery.Recovery[[]*cloudevents.Event]{}, middlewares[0], "the outermost first")
	assert.IsType(t, &logger.Logger[[]*cloudevents.Event]{}, middlewares[1])

	ctx := context.Background()
	out, err := handle(ctx, event(t, map[string]string{"sku": "tshirt"}))
	require.NoError(t, err)
	assert.Empty(t, out, "the reservation is announced by the publisher, not replied")
	left, err := repo.Get(ctx, "tshirt")
	require.NoError(t, err)
	assert.Equal(t, int64(2), left.Units, "the handler took one tshirt from the stock of the app")

	_, err = handle(ctx, event(t, map[string]string{}))
	assert.True(t, errors.IsNotValid(err), "an order with no sku: got %v", err)
	_, err = handle(ctx, event(t, map[string]string{"sku": "nope"}))
	assert.True(t, errors.IsNotFound(err), "an unknown sku: got %v", err)
	_, err = handle(ctx, event(t, map[string]string{"sku": "mug"}))
	require.NoError(t, err)
	_, err = handle(ctx, event(t, map[string]string{"sku": "mug"}))
	assert.True(t, errors.IsConflict(err), "the last mug is gone: got %v", err)
}

const (
	childEnv = "STOCK_EXAMPLE_CHILD"
	// implementationEnv chooses the implementation of reserver.Service
	// (internal.pkg.domain.service.stock.reserver.implementation, local by
	// default).
	implementationEnv = "INTERNAL_PKG_DOMAIN_SERVICE_STOCK_RESERVER_IMPLEMENTATION"
	// hostEnv points the stock client builder generates for the resty
	// reserver at app api (boost.factory.resty.clients.stock.host).
	hostEnv = "BOOST_FACTORY_RESTY_CLIENTS_STOCK_HOST"
)

// child runs test in a new process with env set: config is read once per
// process.
func child(t *testing.T, test string, env ...string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^"+test+"$", "-test.v")
	cmd.Env = append(append(os.Environ(), childEnv+"=1"), env...)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", out)
	assert.Contains(t, string(out), "--- PASS: "+test)
}

// TestWorkerResty boots app worker with the resty reserver, which reserves
// by calling app api: here a server this test runs in its place.
func TestWorkerResty(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method != http.MethodPost:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case r.URL.Path == "/stocks/tshirt/reserve":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"sku":"tshirt","units":2}`))
		case r.URL.Path == "/stocks/mug/reserve":
			w.WriteHeader(http.StatusConflict)
		case r.URL.Path == "/stocks/down/reserve":
			w.WriteHeader(http.StatusBadGateway)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(api.Close)

	child(t, "TestWorkerRestyChild", implementationEnv+"=resty", hostEnv+"="+api.URL)
}

// TestWorkerRestyChild runs in the process TestWorkerResty starts.
func TestWorkerRestyChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestWorkerResty")
	}

	var handle function.Handler[[]*cloudevents.Event]
	app := fx.New(append(worker.Modules(), fx.NopLogger, fx.Populate(&handle))...)
	require.NoError(t, app.Err())

	ctx := context.Background()
	_, err := handle(ctx, event(t, map[string]string{"sku": "tshirt"}))
	assert.NoError(t, err)
	_, err = handle(ctx, event(t, map[string]string{"sku": "mug"}))
	assert.True(t, errors.IsConflict(err), "409 from app api: got %v", err)
	_, err = handle(ctx, event(t, map[string]string{"sku": "nope"}))
	assert.True(t, errors.IsNotFound(err), "404 from app api: got %v", err)
	_, err = handle(ctx, event(t, map[string]string{"sku": "down"}))
	assert.True(t, errors.IsInternal(err), "any other status: got %v", err)
}

// TestWorkerConfig boots the worker with configs/worker.yaml, the file that
// runs it next to app api.
func TestWorkerConfig(t *testing.T) {
	child(t, "TestWorkerConfigChild", "CONF=../../configs/worker.yaml")
}

// TestWorkerConfigChild runs in the process TestWorkerConfig starts.
func TestWorkerConfigChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestWorkerConfig")
	}

	assert.Equal(t, 8081, cehttp.Port(), "the worker listens next to app api on 8080")
	assert.Equal(t, "resty", reserver.ServiceImplementation(), "the worker reserves through app api")

	app := fx.New(append(worker.Modules(), fx.NopLogger)...)
	require.NoError(t, app.Err())
}

// mainEnv makes the process TestMain_ReservesOverHTTP starts run the
// generated main instead of the tests.
const mainEnv = "STOCK_EXAMPLE_MAIN"

// otelEnv turns off the otel plugin the @Server of the cloudevents adapter
// lists (boost.bootstrap.function.adapter.cloudevents.http.plugins.otel.enabled).
const otelEnv = "BOOST_BOOTSTRAP_FUNCTION_ADAPTER_CLOUDEVENTS_HTTP_PLUGINS_OTEL_ENABLED"

// TestMain_ReservesOverHTTP runs the generated main of app worker on a free
// port and posts it an order.created CloudEvent until it reserves the item,
// with the plugins of its @Server and with config turning one off.
func TestMain_ReservesOverHTTP(t *testing.T) {
	t.Run("with the otel plugin", func(t *testing.T) { reservesOverHTTP(t) })
	t.Run("with the otel plugin turned off", func(t *testing.T) { reservesOverHTTP(t, otelEnv+"=false") })
}

func reservesOverHTTP(t *testing.T, env ...string) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := lis.Addr().(*net.TCPAddr).Port
	require.NoError(t, lis.Close())

	out := &syncBuffer{}
	cmd := exec.Command(os.Args[0], "-test.run=^TestMainChild$")
	cmd.Env = append(append(os.Environ(), mainEnv+"=1", "BOOST_BOOTSTRAP_FUNCTION_ADAPTER_CLOUDEVENTS_HTTP_PORT="+strconv.Itoa(port)), env...)
	cmd.Stdout, cmd.Stderr = out, out
	require.NoError(t, cmd.Start())
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	post := func() {
		req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+strconv.Itoa(port)+"/", strings.NewReader(`{"sku":"tshirt"}`))
		require.NoError(t, err)
		req.Header.Set("Ce-Id", "evt-1")
		req.Header.Set("Ce-Source", "test")
		req.Header.Set("Ce-Type", "order.created")
		req.Header.Set("Ce-Specversion", "1.0")
		req.Header.Set("Content-Type", "application/json")
		if res, err := http.DefaultClient.Do(req); err == nil {
			_ = res.Body.Close()
		}
	}

	// the worker listens once it boots: post until it reserves
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	deadline := time.After(10 * time.Second)
	for !strings.Contains(out.String(), "reserved one tshirt for event evt-1") {
		select {
		case <-tick.C:
			post()
		case err := <-exited:
			t.Fatalf("the worker exited (%v) without reserving: %s", err, out.String())
		case <-deadline:
			t.Fatalf("the worker never reserved: %s", out.String())
		}
	}
	assert.Contains(t, out.String(), "listening on :"+strconv.Itoa(port)+"/", "the worker logs the port it listens on")
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
// binary gets, in the process TestMain_ReservesOverHTTP starts.
func TestMainChild(t *testing.T) {
	if os.Getenv(mainEnv) != "1" {
		t.Skip("run by TestMain_ReservesOverHTTP")
	}
	os.Args = os.Args[:1]
	main()
}
