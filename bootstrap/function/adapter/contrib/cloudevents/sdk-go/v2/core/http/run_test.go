package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/model/errors"
)

// TestMain boots boost: the adapter reads its port and path from config.
func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

type customForbidden struct{ msg string }

func (c *customForbidden) Error() string { return c.msg }

func TestErrorStatusCode_RegisteredCustom(t *testing.T) {
	defer errors.ResetRegistry()
	errors.RegisterMatch(func(err error) bool {
		_, ok := err.(*customForbidden)
		return ok
	}, errors.KindForbidden)

	if got := ErrorStatusCode(&customForbidden{"no"}); got != http.StatusForbidden {
		t.Fatalf("ErrorStatusCode = %d, want 403", got)
	}
}

func TestErrorStatusCode_BuiltinNotFound(t *testing.T) {
	if ErrorStatusCode(errors.NotFoundf("x")) != http.StatusNotFound {
		t.Fatal("builtin NotFound not mapped to 404")
	}
}

// withPort is a Plugin that makes the adapter listen on port.
func withPort(port int) Plugin {
	return func(_ context.Context, opts []cehttp.Option) []cehttp.Option {
		return append(opts, cehttp.WithPort(port))
	}
}

// freePort returns a TCP port nothing listens on.
func freePort(t *testing.T) int {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := lis.Addr().(*net.TCPAddr).Port
	require.NoError(t, lis.Close())
	return port
}

// post sends one order.created CloudEvent to port until the adapter
// listens, and returns the status of its reply.
func post(port int) (int, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	tick := time.NewTicker(50 * time.Millisecond)
	defer tick.Stop()
	deadline := time.After(10 * time.Second)
	for {
		req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+strconv.Itoa(port)+"/", strings.NewReader(`{"sku":"tshirt"}`))
		if err != nil {
			return 0, err
		}
		req.Header.Set("Ce-Id", "evt-1")
		req.Header.Set("Ce-Source", "test")
		req.Header.Set("Ce-Type", "order.created")
		req.Header.Set("Ce-Specversion", "1.0")
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err == nil {
			_ = res.Body.Close()
			return res.StatusCode, nil
		}
		select {
		case <-tick.C:
		case <-deadline:
			return 0, err
		}
	}
}

// TestRunContext_LetsTheHandlerFinishAndReturns cancels the run context
// while the handler holds an event (#262): the handler's context is not
// cancelled by the stop, the sender gets the reply and RunContext returns
// nil, so the app stops at once instead of waiting the fx stop timeout.
func TestRunContext_LetsTheHandlerFinishAndReturns(t *testing.T) {
	port := freePort(t)
	held := make(chan context.Context, 1)
	release := make(chan struct{})
	fn := func(ctx context.Context, _ ce.Event) ([]*ce.Event, error) {
		held <- ctx
		<-release
		return nil, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ran := make(chan error, 1)
	go func() { ran <- RunContext[[]*ce.Event](ctx, fn, nil, withPort(port)) }()

	type reply struct {
		status int
		err    error
	}
	replied := make(chan reply, 1)
	go func() {
		status, err := post(port)
		replied <- reply{status, err}
	}()

	var handlerCtx context.Context
	select {
	case handlerCtx = <-held:
	case err := <-ran:
		t.Fatalf("RunContext returned (%v) before the event", err)
	case r := <-replied:
		t.Fatalf("the event was answered (%d, %v) without the handler", r.status, r.err)
	}

	cancel()
	assert.NoError(t, handlerCtx.Err(), "the stop cancelled the event in hand")
	close(release)

	r := <-replied
	require.NoError(t, r.err)
	assert.Equal(t, http.StatusOK, r.status, "the sender gets the reply of the event in hand")
	select {
	case err := <-ran:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("RunContext did not return after the run context was cancelled")
	}
}

// TestRunContext_FailsOnABadProtocolOption returns the error of an option
// the protocol rejects, instead of logging it and failing on a nil protocol.
func TestRunContext_FailsOnABadProtocolOption(t *testing.T) {
	fn := func(context.Context, ce.Event) ([]*ce.Event, error) { return nil, nil }

	err := RunContext[[]*ce.Event](context.Background(), fn, nil, withPort(-1))
	assert.ErrorContains(t, err, "invalid port: -1")

	err = Run[[]*ce.Event](fn, nil, withPort(-1))
	assert.ErrorContains(t, err, "invalid port: -1", "Run is RunContext with a context never cancelled")
}

// TestRunContext_FailsOnABadClientOption returns the error of a client
// option that fails.
func TestRunContext_FailsOnABadClientOption(t *testing.T) {
	fn := func(context.Context, ce.Event) ([]*ce.Event, error) { return nil, nil }
	bad := func(any) error { return errors.New("bad client option") }

	err := RunContext[[]*ce.Event](context.Background(), fn, []client.Option{bad}, withPort(freePort(t)))
	assert.ErrorContains(t, err, "bad client option")
}

// TestNew_RunsOnTheCommandContext executes the cloudevents_http command with
// a context already done (#262): it runs on the command's context, so it
// returns instead of serving forever.
func TestNew_RunsOnTheCommandContext(t *testing.T) {
	fn := func(context.Context, ce.Event) ([]*ce.Event, error) { return nil, nil }
	cmd := New[[]*ce.Event](nil, withPort(freePort(t)))(fn)
	assert.Equal(t, "cloudevents_http", cmd.Use)
	cmd.SetArgs([]string{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ran := make(chan error, 1)
	go func() { ran <- cmd.ExecuteContext(ctx) }()
	select {
	case err := <-ran:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("the command did not return on a context already done")
	}
}
