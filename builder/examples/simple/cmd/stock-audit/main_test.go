package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/job"
	stockaudit "github.com/xgodev/boost/builder/examples/simple/gen/builder/stock-audit"
	"github.com/xgodev/boost/builder/examples/simple/stock"
	"go.uber.org/fx"
)

// TestMain boots boost: the stock and the audit read their config.
func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// TestAudit builds app stock-audit from the generated Modules() and runs the
// job handler it provides, as the job module does once on start.
func TestAudit(t *testing.T) {
	var (
		run      job.Handler
		stockSvc stock.Service
	)
	app := fx.New(append(stockaudit.Modules(), fx.NopLogger, fx.Populate(&run, &stockSvc))...)
	require.NoError(t, app.Err())

	assert.Contains(t, reflect.TypeOf(stockSvc).Elem().PkgPath(), "/gen/builder/middleware/",
		"the auditor gets the stock wrapped in the middlewares of @Middleware")
	require.NoError(t, run(context.Background()), "the memory stock answers every item")
}

const (
	// mainEnv makes the process the TestMain_* tests start run the generated
	// main instead of the tests.
	mainEnv = "BUILDER_EXAMPLE_MAIN"
	// implementationEnv chooses the implementation of stock.Service
	// (builder.examples.simple.stock.implementation).
	implementationEnv = "BUILDER_EXAMPLES_SIMPLE_STOCK_IMPLEMENTATION"
	// warehouseHostEnv points the warehouse client of stock/remote at a
	// server (boost.factory.resty.clients.warehouse.host).
	warehouseHostEnv = "BOOST_FACTORY_RESTY_CLIENTS_WAREHOUSE_HOST"
	// requestTimeoutEnv sets how long the warehouse client waits for an
	// answer (boost.factory.resty.clients.warehouse.requestTimeout).
	requestTimeoutEnv = "BOOST_FACTORY_RESTY_CLIENTS_WAREHOUSE_REQUEST__TIMEOUT"
)

// TestMain_ExitsZeroWhenTheJobSucceeds runs the generated main of app
// stock-audit on the memory stock: the audit returns nil and the process
// exits 0 by itself, nothing stops it.
func TestMain_ExitsZeroWhenTheJobSucceeds(t *testing.T) {
	out, err := runMain(t, 10*time.Second)

	require.NoError(t, err, "the audit did not exit 0: %s", out)
	assert.Contains(t, out, "sold-out is out of stock")
	assert.Contains(t, out, "audited 3 items, 1 out of stock")
}

// TestMain_ExitsOneWhenTheJobFails runs the generated main of app stock-audit
// on the remote stock with a warehouse that answers 500: the audit returns
// the error, which the app logs before exiting 1.
func TestMain_ExitsOneWhenTheJobFails(t *testing.T) {
	warehouse := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(warehouse.Close)

	out, err := runMain(t, 10*time.Second, implementationEnv+"=remote", warehouseHostEnv+"="+warehouse.URL)

	var exit *exec.ExitError
	require.ErrorAs(t, err, &exit, "the audit must fail: %s", out)
	assert.Equal(t, 1, exit.ExitCode())
	assert.Contains(t, out, "job failed")
	assert.Contains(t, out, "the warehouse answered 500")
}

// TestMain_SIGTERMCancelsTheJob sends SIGTERM to the generated main of app
// stock-audit while the audit waits on a warehouse that never answers: the
// job context is cancelled, the request returns, and the process exits 1
// right away instead of waiting the fx stop timeout, saying the job did not
// finish.
func TestMain_SIGTERMCancelsTheJob(t *testing.T) {
	asked := make(chan struct{}, 1)
	warehouse := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		asked <- struct{}{}
		<-r.Context().Done()
	}))
	t.Cleanup(warehouse.Close)

	out := &syncBuffer{}
	cmd := exec.Command(os.Args[0], "-test.run=^TestMainChild$")
	cmd.Env = append(os.Environ(), implementationEnv+"=remote", warehouseHostEnv+"="+warehouse.URL,
		requestTimeoutEnv+"=1m", mainEnv+"=1")
	cmd.Stdout, cmd.Stderr = out, out
	require.NoError(t, cmd.Start())
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	select {
	case <-asked:
	case err := <-exited:
		t.Fatalf("the audit exited before asking the warehouse: %v: %s", err, out.String())
	case <-time.After(10 * time.Second):
		t.Fatalf("the audit did not ask the warehouse: %s", out.String())
	}

	require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))
	select {
	case err := <-exited:
		var exit *exec.ExitError
		require.ErrorAs(t, err, &exit, "an audit cut short must not exit 0: %s", out.String())
		assert.Equal(t, 1, exit.ExitCode())
		assert.Contains(t, out.String(), "job stopped before finishing")
		assert.Contains(t, out.String(), "context canceled")
	case <-time.After(5 * time.Second):
		t.Fatalf("the audit did not stop within 5s of SIGTERM: %s", out.String())
	}
}

// runMain runs the generated main of app stock-audit with env and returns
// its output and exit, failing the test if it runs past timeout.
func runMain(t *testing.T, timeout time.Duration, env ...string) (string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestMainChild$")
	cmd.Env = append(append(os.Environ(), env...), mainEnv+"=1")
	out, err := cmd.CombinedOutput()
	require.NoError(t, ctx.Err(), "the audit did not exit by itself: %s", out)
	return string(out), err
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
