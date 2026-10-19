package server

import (
	"context"
	"errors"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/storj.io/drpc/v1/client"
	"github.com/xgodev/boost/wrapper/config"
	"storj.io/drpc"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// wait bounds every channel wait in these tests.
const wait = 5 * time.Second

// text is a drpc encoding for *string messages.
type text struct{}

func (text) Marshal(msg drpc.Message) ([]byte, error)  { return []byte(*msg.(*string)), nil }
func (text) Unmarshal(b []byte, msg drpc.Message) error { *msg.(*string) = string(b); return nil }

// echo replies with the request prefixed by "echo: ".
type echo struct{}

func (echo) Echo(_ context.Context, in *string) (*string, error) {
	out := "echo: " + *in
	return &out, nil
}

// echoDesc describes the single /echo/Echo rpc.
type echoDesc struct{}

func (echoDesc) NumMethods() int { return 1 }

func (echoDesc) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	if n != 0 {
		return "", nil, nil, nil, false
	}
	return "/echo/Echo", text{},
		func(srv interface{}, ctx context.Context, in1, _ interface{}) (drpc.Message, error) {
			return srv.(echo).Echo(ctx, in1.(*string))
		}, echo.Echo, true
}

// serve builds a server with plugins, runs Serve in the background and
// returns its port once it listens.
func serve(t *testing.T, plugins ...Plugin) (*Server, int) {
	t.Helper()

	opts, err := NewOptions()
	require.NoError(t, err)
	opts.Port = 0

	srv, err := NewServerWithOptions(context.Background(), opts, plugins...)
	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		srv.Serve(context.Background())
		close(done)
	}()
	t.Cleanup(func() {
		srv.Shutdown(context.Background())
		select {
		case <-done:
		case <-time.After(wait):
			t.Error("Serve did not return after Shutdown")
		}
	})

	select {
	case <-srv.listening:
	case <-time.After(wait):
		t.Fatal("server did not listen")
	}
	return srv, srv.addr.(*net.TCPAddr).Port
}

func dial(t *testing.T, port int) drpc.Conn {
	t.Helper()
	conn, err := client.NewClientConnWithOptions(context.Background(), &client.Options{Host: "127.0.0.1", Port: port, Timeout: wait})
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

func TestServeAnswersRegisteredRPC(t *testing.T) {
	register := func(_ context.Context, s *Server) error { return s.Register(echo{}, echoDesc{}) }
	_, port := serve(t, register)

	in, out := "ping", ""
	err := dial(t, port).Invoke(context.Background(), "/echo/Echo", text{}, &in, &out)

	require.NoError(t, err)
	assert.Equal(t, "echo: ping", out)
}

func TestServeRejectsUnknownRPC(t *testing.T) {
	_, port := serve(t)

	in, out := "ping", ""
	err := dial(t, port).Invoke(context.Background(), "/echo/Echo", text{}, &in, &out)

	assert.Error(t, err)
}

func TestFailingPluginStopsTheBuild(t *testing.T) {
	boom := errors.New("boom")
	ran := false
	failing := func(context.Context, *Server) error { return boom }
	later := func(context.Context, *Server) error { ran = true; return nil }

	srv, err := NewServerWithOptions(context.Background(), &Options{}, failing, later)

	assert.Nil(t, srv)
	assert.ErrorIs(t, err, boom)
	assert.False(t, ran)
}

func TestShutdownBeforeServeMakesServeReturn(t *testing.T) {
	srv, err := NewServerWithOptions(context.Background(), &Options{})
	require.NoError(t, err)
	srv.Shutdown(context.Background())

	done := make(chan struct{})
	go func() {
		srv.Serve(context.Background())
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(wait):
		t.Fatal("Serve did not return")
	}
}

func TestShutdownClosesTheListener(t *testing.T) {
	srv, port := serve(t)
	srv.Shutdown(context.Background())

	_, err := client.NewClientConnWithOptions(context.Background(), &client.Options{Host: "127.0.0.1", Port: port, Timeout: time.Second})

	assert.Error(t, err)
}

func TestNewOptionsWithPathReadsTheEnv(t *testing.T) {
	t.Cleanup(config.Load)
	t.Setenv("APP_DRPC_PORT", "7001")
	t.Setenv("APP_DRPC_COLLECT__STATS", "true")

	opts, err := NewOptionsWithPath("app.drpc")

	require.NoError(t, err)
	assert.Equal(t, 7001, opts.Port)
	assert.True(t, opts.CollectStats)
}

func TestNewOptionsDefaults(t *testing.T) {
	opts, err := NewOptions()

	require.NoError(t, err)
	assert.Equal(t, 9091, opts.Port)
	assert.False(t, opts.CollectStats)
}
