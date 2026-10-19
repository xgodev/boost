package client

import (
	"context"
	stderrors "errors"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func TestNewClientUsesTheConfigDefaults(t *testing.T) {
	client, err := NewClient(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "fasthttp-client", client.Name)
	assert.Equal(t, 512, client.MaxConnsPerHost)
	assert.Equal(t, 10*time.Second, client.MaxIdleConnDuration)
	assert.Equal(t, 52428800, client.MaxResponseBodySize)
	assert.Equal(t, 5, client.MaxIdemponentCallAttempts)
}

func TestNewClientWithConfigPathReadsTheGivenPath(t *testing.T) {
	t.Cleanup(config.Load)
	t.Setenv("APP_HTTP_NAME", "orders")
	t.Setenv("APP_HTTP_READ__TIMEOUT", "3s")

	client, err := NewClientWithConfigPath(context.Background(), "app.http")
	require.NoError(t, err)

	assert.Equal(t, "orders", client.Name)
	assert.Equal(t, 3*time.Second, client.ReadTimeout)
	assert.Equal(t, 512, client.MaxConnsPerHost, "keys not under the path keep the defaults")
}

func TestPluginsChangeTheClient(t *testing.T) {
	client, err := NewClient(context.Background(), func(_ context.Context, c *fasthttp.Client) error {
		c.NoDefaultUserAgentHeader = true
		return nil
	})
	require.NoError(t, err)

	assert.True(t, client.NoDefaultUserAgentHeader)
}

func TestAFailingPluginStopsTheBuild(t *testing.T) {
	boom := stderrors.New("boom")

	_, err := NewClient(context.Background(), func(context.Context, *fasthttp.Client) error { return boom })

	assert.ErrorIs(t, err, boom)
}

func TestTheClientSendsRequests(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	require.NoError(t, err)
	srv := &fasthttp.Server{Handler: func(c *fasthttp.RequestCtx) {
		c.SetBodyString(string(c.UserAgent()))
	}}
	go func() { _ = srv.Serve(ln) }()
	t.Cleanup(func() { _ = srv.Shutdown() })

	client, err := NewClient(context.Background())
	require.NoError(t, err)

	status, body, err := client.Get(nil, "http://"+ln.Addr().String()+"/")
	require.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, status)
	assert.Equal(t, "fasthttp-client", string(body), "the name is the User-Agent")
}
