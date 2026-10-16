package multiserver

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	server "github.com/xgodev/boost/extra/multiserver"
	fxcontext "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

const wait = 5 * time.Second

// TestMain spends the core context module's once-per-process registration,
// so every app below supplies its own context and none depends on test order.
func TestMain(m *testing.M) {
	fxcontext.Module()
	os.Exit(m.Run())
}

// fakeServer behaves like a real listener: Serve blocks until Shutdown
// releases it (or, when stuck, until the test ends).
type fakeServer struct {
	stuck     bool
	served    chan struct{} // closed when Serve starts
	serveDone chan struct{} // closed when Serve returns
	release   chan struct{} // closed to unblock Serve

	mu          sync.Mutex
	serveCtx    context.Context
	shutdownCtx context.Context
	releaseOnce sync.Once
}

func newFakeServer(t *testing.T, stuck bool) *fakeServer {
	s := &fakeServer{
		stuck:     stuck,
		served:    make(chan struct{}),
		serveDone: make(chan struct{}),
		release:   make(chan struct{}),
	}
	t.Cleanup(s.unblock)
	return s
}

func (s *fakeServer) unblock() { s.releaseOnce.Do(func() { close(s.release) }) }

func (s *fakeServer) Serve(ctx context.Context) {
	s.mu.Lock()
	s.serveCtx = ctx
	s.mu.Unlock()
	close(s.served)
	defer close(s.serveDone)
	<-s.release
}

func (s *fakeServer) Shutdown(ctx context.Context) {
	s.mu.Lock()
	s.shutdownCtx = ctx
	s.mu.Unlock()
	if !s.stuck {
		s.unblock()
	}
}

// newApp builds the module the way fxfact.NewApp does and fails if building
// blocks: construction must not serve, or lifecycle hooks never run.
func newApp(t *testing.T, srvs ...server.Server) *fx.App {
	t.Helper()
	once = sync.Once{} // Module() registers once per process

	built := make(chan *fx.App, 1)
	go func() {
		built <- fx.New(
			fx.NopLogger,
			Module(),
			fx.Provide(
				context.Background,
				fx.Annotate(
					func() []server.Server { return srvs },
					fx.ResultTags(`group:"extra.multiserver.servers,flatten"`),
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
		t.Fatal("fx.New did not return: servers are served during construction, so fx lifecycle hooks never run")
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

func TestModuleServesOnStartAndDrainsOnStop(t *testing.T) {
	srv := newFakeServer(t, false)
	app := newApp(t, srv)

	select {
	case <-srv.served:
		t.Fatal("Serve was called before app.Start")
	default:
	}

	startApp(t, app)
	waitClosed(t, srv.served, "Serve was not called after app.Start")

	stopCtx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	stopDeadline, _ := stopCtx.Deadline()
	if err := app.Stop(stopCtx); err != nil {
		t.Fatalf("app.Stop: %v", err)
	}

	select {
	case <-srv.serveDone:
	default:
		t.Fatal("app.Stop returned before Serve returned: in-flight work is not drained")
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.shutdownCtx == nil {
		t.Fatal("Shutdown was not called on Stop")
	}
	if got, ok := srv.shutdownCtx.Deadline(); !ok || !got.Equal(stopDeadline) {
		t.Fatalf("Shutdown context deadline = %v (set=%v), want the stop context deadline %v", got, ok, stopDeadline)
	}
	if srv.serveCtx.Err() == nil {
		t.Fatal("the context passed to Serve is still alive after Stop")
	}
}

func TestModuleStopIsBoundedByStopContext(t *testing.T) {
	srv := newFakeServer(t, true)
	app := newApp(t, srv)
	startApp(t, app)
	waitClosed(t, srv.served, "Serve was not called after app.Start")

	stopCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	stopped := make(chan error, 1)
	go func() { stopped <- app.Stop(stopCtx) }()

	select {
	case err := <-stopped:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("app.Stop error = %v, want context.DeadlineExceeded for a server that does not drain", err)
		}
	case <-time.After(wait):
		t.Fatal("app.Stop ignored the stop context deadline while a server did not drain")
	}

	// extra/multiserver keeps package-level state: let Serve return before
	// the next test serves again.
	srv.unblock()
	waitClosed(t, srv.serveDone, "Serve did not return after being released")
}

func TestModuleShutsDownWhenServersStopOnTheirOwn(t *testing.T) {
	srv := newFakeServer(t, false)
	app := newApp(t, srv)
	startApp(t, app)
	waitClosed(t, srv.served, "Serve was not called after app.Start")

	srv.unblock() // the server dies without Stop being requested

	select {
	case sig := <-app.Wait():
		if sig.ExitCode != 1 {
			t.Fatalf("shutdown exit code = %d, want 1 for servers that stopped unexpectedly", sig.ExitCode)
		}
	case <-time.After(wait):
		t.Fatal("every server stopped but the app was not asked to shut down")
	}

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		t.Fatalf("app.Stop: %v", err)
	}
}
