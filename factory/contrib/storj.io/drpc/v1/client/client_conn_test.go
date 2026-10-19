package client

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
	"github.com/xgodev/boost/wrapper/config"
	"storj.io/drpc/drpcconn"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// listen returns the port of a listener that accepts and holds connections.
func listen(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			t.Cleanup(func() { _ = c.Close() })
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestNewClientConnWithOptionsRunsPlugins(t *testing.T) {
	port := listen(t)
	var got *drpcconn.Conn
	plugin := func(_ context.Context, c *drpcconn.Conn) error { got = c; return nil }

	conn, err := NewClientConnWithOptions(context.Background(), &Options{Host: "127.0.0.1", Port: port, Timeout: time.Second}, plugin)

	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	assert.Same(t, conn, got)
}

func TestFailingPluginClosesTheConn(t *testing.T) {
	port := listen(t)
	boom := errors.New("boom")
	var got *drpcconn.Conn
	failing := func(_ context.Context, c *drpcconn.Conn) error { got = c; return boom }

	conn, err := NewClientConnWithOptions(context.Background(), &Options{Host: "127.0.0.1", Port: port, Timeout: time.Second}, failing)

	assert.Nil(t, conn)
	assert.ErrorIs(t, err, boom)
	select {
	case <-got.Closed():
	case <-time.After(time.Second):
		t.Fatal("conn not closed")
	}
}

func TestNewClientConnWithOptionsFailsWithoutServer(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	require.NoError(t, ln.Close())

	_, err = NewClientConnWithOptions(context.Background(), &Options{Host: "127.0.0.1", Port: port, Timeout: time.Second})

	assert.Error(t, err)
}

func TestNewOptionsDefaults(t *testing.T) {
	opts, err := NewOptions()

	require.NoError(t, err)
	assert.Equal(t, "localhost", opts.Host)
	assert.Equal(t, 9091, opts.Port)
	assert.Equal(t, 10*time.Second, opts.Timeout)
}

func TestNewOptionsWithPathReadsTheEnv(t *testing.T) {
	t.Cleanup(config.Load)
	t.Setenv("APP_DRPC_HOST", "orders")
	t.Setenv("APP_DRPC_PORT", "7001")
	t.Setenv("APP_DRPC_TIMEOUT", "3s")

	opts, err := NewOptionsWithPath("app.drpc")

	require.NoError(t, err)
	assert.Equal(t, "orders", opts.Host)
	assert.Equal(t, 7001, opts.Port)
	assert.Equal(t, 3*time.Second, opts.Timeout)
}
