package function

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/extra/middleware"
	fxcontext "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

const wait = 5 * time.Second

// TestMain boots boost (the function reads its config) and spends the core
// context module's once-per-process registration, so every app below
// supplies its own context and none depends on test order.
func TestMain(m *testing.M) {
	boost.Start()
	fxcontext.Module()
	os.Exit(m.Run())
}

// adapter is a function adapter command whose run is observable: it reports
// when it starts and returns, and runs body with the command context and,
// in handler, the handler the function gave the command.
type adapter struct {
	started  chan struct{}
	returned chan struct{}
	body     func(ctx context.Context) error
	handler  function.Handler[*event.Event]
}

func newAdapter(body func(ctx context.Context) error) *adapter {
	return &adapter{started: make(chan struct{}), returned: make(chan struct{}), body: body}
}

func (a *adapter) cmd(h function.Handler[*event.Event]) *co.Command {
	a.handler = h
	return &co.Command{
		Use: "fake",
		RunE: func(cmd *co.Command, _ []string) error {
			close(a.started)
			defer close(a.returned)
			return a.body(cmd.Context())
		},
	}
}

// newApp builds the module the way fxfact.NewApp does and fails if building
// blocks: construction must not run the function, or lifecycle hooks never run.
func newApp(t *testing.T, a *adapter) *fx.App {
	t.Helper()
	return newAppWith(t, a, nil)
}

// newAppWith is newApp with the middlewares given to Module and more options.
func newAppWith(t *testing.T, a *adapter, m []middleware.AnyErrorMiddleware[*event.Event], opts ...fx.Option) *fx.App {
	t.Helper()
	once = sync.Once{} // Module() registers once per process

	args := os.Args
	os.Args = []string{"func", "fake"}
	t.Cleanup(func() { os.Args = args })

	built := make(chan *fx.App, 1)
	go func() {
		built <- fx.New(
			fx.NopLogger,
			Module[*event.Event](m),
			fx.Options(opts...),
			fx.Provide(
				context.Background,
				func() function.Handler[*event.Event] {
					return func(context.Context, cloudevents.Event) (*event.Event, error) { return nil, nil }
				},
				fx.Annotate(
					func() function.CmdFunc[*event.Event] { return a.cmd },
					fx.ResultTags(`group:"`+BSFunctionAdaptersGroupKey+`"`),
				),
			),
		)
	}()

	select {
	case app := <-built:
		if err := app.Err(); err != nil {
			t.Fatalf("fx.New: %v", err)
		}
		return app
	case <-time.After(wait):
		t.Fatal("fx.New did not return: the function runs during construction, so fx lifecycle hooks never run")
		return nil
	}
}

func waitClosed(t *testing.T, ch <-chan struct{}, msg string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(wait):
		t.Fatal(msg)
	}
}

func startApp(t *testing.T, app *fx.App) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err := app.Start(ctx); err != nil {
		t.Fatalf("app.Start: %v", err)
	}
}

func stopApp(t *testing.T, app *fx.App) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	return app.Stop(ctx)
}

func TestModuleRunsOnStartAndStopsOnStop(t *testing.T) {
	a := newAdapter(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})
	app := newApp(t, a)

	select {
	case <-a.started:
		t.Fatal("the function ran before app.Start")
	default:
	}

	startApp(t, app)
	waitClosed(t, a.started, "the function did not run after app.Start")

	if err := stopApp(t, app); err != nil {
		t.Fatalf("app.Stop: %v", err)
	}
	select {
	case <-a.returned:
	default:
		t.Fatal("app.Stop returned before the function returned: in-flight work is not drained")
	}
}

func TestModuleShutsDownWhenFunctionEnds(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{name: "success", err: nil, wantCode: 0},
		{name: "failure", err: errors.New("boom"), wantCode: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newAdapter(func(context.Context) error { return tt.err })
			app := newApp(t, a)
			startApp(t, app)

			select {
			case sig := <-app.Wait():
				if sig.ExitCode != tt.wantCode {
					t.Fatalf("shutdown exit code = %d, want %d", sig.ExitCode, tt.wantCode)
				}
			case <-time.After(wait):
				t.Fatal("the function ended but the app was not asked to shut down")
			}
			if err := stopApp(t, app); err != nil {
				t.Fatalf("app.Stop: %v", err)
			}
		})
	}
}

// recorder is a middleware that appends its name to calls and runs the next.
type recorder struct {
	name  string
	calls *[]string
}

func (r recorder) Exec(ctx *middleware.AnyErrorContext[*event.Event], exec middleware.AnyErrorExecFunc[*event.Event], fallback middleware.AnyErrorReturnFunc[*event.Event]) (*event.Event, error) {
	*r.calls = append(*r.calls, r.name)
	return ctx.Next(exec, fallback)
}

func TestModuleRunsTheGivenThenTheNamedMiddlewares(t *testing.T) {
	var calls []string
	var a *adapter
	a = newAdapter(func(ctx context.Context) error {
		_, err := a.handler(ctx, cloudevents.NewEvent())
		return err
	})
	app := newAppWith(t, a,
		[]middleware.AnyErrorMiddleware[*event.Event]{recorder{"given", &calls}},
		fx.Provide(fx.Annotate(
			func() []middleware.AnyErrorMiddleware[*event.Event] {
				return []middleware.AnyErrorMiddleware[*event.Event]{recorder{"first", &calls}, recorder{"second", &calls}}
			},
			fx.ResultTags(`name:"`+MiddlewaresName+`"`),
		)),
	)
	startApp(t, app)
	waitClosed(t, a.returned, "the function did not run")
	if err := stopApp(t, app); err != nil {
		t.Fatalf("app.Stop: %v", err)
	}

	if got, want := strings.Join(calls, ","), "given,first,second"; got != want {
		t.Fatalf("middlewares ran as %q, want %q", got, want)
	}
}
