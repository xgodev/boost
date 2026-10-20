package job

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/job"
	fxcontext "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

const wait = 5 * time.Second

// TestMain boots boost (the module logs) and spends the core context
// module's once-per-process registration, so every app below supplies its
// own context.
func TestMain(m *testing.M) {
	boost.Start()
	fxcontext.Module()
	os.Exit(m.Run())
}

// newApp builds the module with the handler; building must not run it.
func newApp(t *testing.T, h job.Handler) *fx.App {
	t.Helper()
	once = sync.Once{} // Module() registers once per process

	app := fx.New(fx.NopLogger, Module(), fx.Provide(context.Background, func() job.Handler { return h }))
	require.NoError(t, app.Err())
	return app
}

func start(t *testing.T, app *fx.App) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	require.NoError(t, app.Start(ctx))
}

// shutdown waits for the signal the job sends when its handler returns.
func shutdown(t *testing.T, app *fx.App) fx.ShutdownSignal {
	t.Helper()
	select {
	case sig := <-app.Wait():
		return sig
	case <-time.After(wait):
		t.Fatal("the job did not shut the app down")
		return fx.ShutdownSignal{}
	}
}

func stop(app *fx.App, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return app.Stop(ctx)
}

func TestModule_RunsOnceOnStartAndExitsZero(t *testing.T) {
	runs := 0
	app := newApp(t, func(context.Context) error {
		runs++
		return nil
	})
	assert.Zero(t, runs, "the handler ran while building the app")

	start(t, app)
	sig := shutdown(t, app)
	assert.Equal(t, 0, sig.ExitCode)
	require.NoError(t, stop(app, wait))
	assert.Equal(t, 1, runs)
}

func TestModule_ErrorExitsOne(t *testing.T) {
	app := newApp(t, func(context.Context) error { return errors.New("no rows") })

	start(t, app)
	sig := shutdown(t, app)
	assert.Equal(t, 1, sig.ExitCode)
	require.NoError(t, stop(app, wait), "the error is reported once, by the exit code and its log line")
}

func TestModule_StopCancelsTheRunAndFails(t *testing.T) {
	cancelled := make(chan struct{})
	app := newApp(t, func(ctx context.Context) error {
		<-ctx.Done()
		close(cancelled)
		return ctx.Err()
	})

	start(t, app)
	err := stop(app, wait)
	require.Error(t, err, "a job stopped before it finished must not exit 0")
	assert.ErrorIs(t, err, context.Canceled)
	select {
	case <-cancelled:
	default:
		t.Fatal("stop did not cancel the context of the handler")
	}
}

func TestModule_StopAfterACleanReturnSucceeds(t *testing.T) {
	app := newApp(t, func(ctx context.Context) error {
		<-ctx.Done()
		return nil // finished its work before giving up
	})

	start(t, app)
	require.NoError(t, stop(app, wait))
}

func TestModule_StopGivesUpAtTheStopTimeout(t *testing.T) {
	release := make(chan struct{})
	t.Cleanup(func() { close(release) })
	app := newApp(t, func(context.Context) error {
		<-release // ignores its context
		return nil
	})

	start(t, app)
	began := time.Now()
	err := stop(app, 100*time.Millisecond)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Less(t, time.Since(began), wait)
}
