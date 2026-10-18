package server

import (
	"context"
	"fmt"
	"net"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
)

func loadConfig(t *testing.T) {
	t.Helper()
	config.Set(koanf.New())
	config.Load()
}

type echoIn struct {
	Text string `json:"text"`
}

type echoOut struct {
	Text string `json:"text"`
}

func addEcho(srv *Server) {
	mcp.AddTool(srv.Instance(), &mcp.Tool{Name: "echo"},
		func(_ context.Context, _ *mcp.CallToolRequest, in echoIn) (*mcp.CallToolResult, echoOut, error) {
			return nil, echoOut{Text: in.Text}, nil
		})
}

func callEcho(t *testing.T, ctx context.Context, transport mcp.Transport) {
	t.Helper()
	cli := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	session, err := cli.Connect(ctx, transport, nil)
	require.NoError(t, err)
	defer session.Close()

	res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "echo", Arguments: map[string]any{"text": "hi"}})
	require.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Equal(t, map[string]any{"text": "hi"}, res.StructuredContent)
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func httpOptions(port int) *Options {
	return &Options{Name: "test", Version: "1.0.0", Transport: TransportHTTP, Port: port, Path: "/mcp"}
}

func TestNewOptions_Defaults(t *testing.T) {
	loadConfig(t)

	o, err := NewOptions()
	require.NoError(t, err)
	assert.Equal(t, boost.ApplicationName(), o.Name)
	assert.Equal(t, "1.0.0", o.Version)
	assert.Equal(t, TransportHTTP, o.Transport)
	assert.Equal(t, 8090, o.Port)
	assert.Equal(t, "/mcp", o.Path)
	assert.False(t, o.Stateless)
	assert.False(t, o.JSONResponse)
	assert.Zero(t, o.SessionTimeout)
	assert.Zero(t, o.KeepAlive)
}

func TestNewOptionsWithPath_Overrides(t *testing.T) {
	t.Setenv("BOOST_FACTORY_MCP_SERVER_ADMIN_PORT", "9999")
	t.Setenv("BOOST_FACTORY_MCP_SERVER_ADMIN_JSON__RESPONSE", "true")
	t.Setenv("BOOST_FACTORY_MCP_SERVER_ADMIN_SESSION__TIMEOUT", "30s")
	ConfigAdd(root + ".admin")
	loadConfig(t)

	o, err := NewOptionsWithPath(root + ".admin")
	require.NoError(t, err)
	assert.Equal(t, 9999, o.Port)
	assert.True(t, o.JSONResponse)
	assert.Equal(t, 30*time.Second, o.SessionTimeout)
	assert.Equal(t, "/mcp", o.Path)
}

func TestNewServer_FromConfig(t *testing.T) {
	loadConfig(t)
	ctx := context.Background()

	srv, err := NewServer(ctx)
	require.NoError(t, err)
	assert.NotNil(t, srv.Instance())
	assert.Equal(t, 8090, srv.Options().Port)

	srv, err = NewServerWithConfigPath(ctx, root)
	require.NoError(t, err)
	assert.Equal(t, TransportHTTP, srv.Options().Transport)
}

func TestNewServerWithOptions_UnknownTransport(t *testing.T) {
	o := httpOptions(0)
	o.Transport = "SSE"

	_, err := NewServerWithOptions(context.Background(), o)
	assert.True(t, errors.IsNotValid(err), "got %v", err)
}

func TestNewServerWithOptions_TransportIsCaseInsensitive(t *testing.T) {
	o := httpOptions(0)
	o.Transport = "stdio"

	srv, err := NewServerWithOptions(context.Background(), o)
	require.NoError(t, err)
	assert.Equal(t, TransportSTDIO, srv.Options().Transport)
}

func TestNewServerWithOptions_PluginRunsAndFails(t *testing.T) {
	ctx := context.Background()
	var got *Server
	ok := func(_ context.Context, s *Server) error { got = s; return nil }

	srv, err := NewServerWithOptions(ctx, httpOptions(0), ok)
	require.NoError(t, err)
	assert.Same(t, srv, got)

	boom := fmt.Errorf("boom")
	_, err = NewServerWithOptions(ctx, httpOptions(0), func(context.Context, *Server) error { return boom })
	assert.ErrorIs(t, err, boom)
}

func TestHandler_ServesToolsOverStreamableHTTP(t *testing.T) {
	ctx := context.Background()
	srv, err := NewServerWithOptions(ctx, httpOptions(0))
	require.NoError(t, err)
	addEcho(srv)

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	callEcho(t, ctx, &mcp.StreamableClientTransport{Endpoint: ts.URL, MaxRetries: -1})
}

func TestServe_HTTP_ServesAndShutsDown(t *testing.T) {
	ctx := context.Background()
	port := freePort(t)
	srv, err := NewServerWithOptions(ctx, httpOptions(port))
	require.NoError(t, err)
	addEcho(srv)

	done := make(chan struct{})
	go func() { srv.Serve(ctx); close(done) }()

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/mcp", port)
	require.Eventually(t, func() bool {
		c, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			c.Close()
		}
		return err == nil
	}, 2*time.Second, 10*time.Millisecond)

	callEcho(t, ctx, &mcp.StreamableClientTransport{Endpoint: endpoint, MaxRetries: -1})

	srv.Shutdown(ctx)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Serve did not return after Shutdown")
	}
}

func TestServe_HTTP_ListenErrorReturns(t *testing.T) {
	ctx := context.Background()
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()

	srv, err := NewServerWithOptions(ctx, httpOptions(l.Addr().(*net.TCPAddr).Port))
	require.NoError(t, err)

	done := make(chan struct{})
	go func() { srv.Serve(ctx); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Serve should return when the port is taken")
	}
	srv.Shutdown(ctx) // no-op once Serve gave up
}

func TestServe_STDIO_ServesAndShutsDown(t *testing.T) {
	ctx := context.Background()
	o := httpOptions(0)
	o.Transport = TransportSTDIO
	srv, err := NewServerWithOptions(ctx, o)
	require.NoError(t, err)
	addEcho(srv)

	serverSide, clientSide := mcp.NewInMemoryTransports()
	srv.transport = serverSide

	done := make(chan struct{})
	go func() { srv.Serve(ctx); close(done) }()

	callEcho(t, ctx, clientSide)

	srv.Shutdown(ctx)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Serve did not return after Shutdown")
	}
}

func TestShutdown_BeforeServeIsNoop(t *testing.T) {
	srv, err := NewServerWithOptions(context.Background(), httpOptions(0))
	require.NoError(t, err)
	srv.Shutdown(context.Background())
}

func TestNewServer_ConfigErrorIsReturned(t *testing.T) {
	t.Setenv("BOOST_FACTORY_MCP_SERVER_PORT", "not-a-port")
	loadConfig(t)
	ctx := context.Background()

	_, err := NewServer(ctx)
	assert.Error(t, err)
	_, err = NewServerWithConfigPath(ctx, root)
	assert.Error(t, err)
}

type failingTransport struct{}

func (failingTransport) Connect(context.Context) (mcp.Connection, error) {
	return nil, fmt.Errorf("no stdin")
}

func TestServe_STDIO_TransportErrorReturns(t *testing.T) {
	o := httpOptions(0)
	o.Transport = TransportSTDIO
	srv, err := NewServerWithOptions(context.Background(), o)
	require.NoError(t, err)
	srv.transport = failingTransport{}

	srv.Serve(context.Background()) // returns instead of blocking
}

func TestShutdown_HTTP_ExpiredContextWithOpenSession(t *testing.T) {
	ctx := context.Background()
	port := freePort(t)
	srv, err := NewServerWithOptions(ctx, httpOptions(port))
	require.NoError(t, err)

	done := make(chan struct{})
	go func() { srv.Serve(ctx); close(done) }()
	require.Eventually(t, func() bool {
		srv.mu.Lock()
		defer srv.mu.Unlock()
		return srv.http != nil
	}, 2*time.Second, 10*time.Millisecond)

	cli := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	session, err := cli.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint: fmt.Sprintf("http://127.0.0.1:%d/mcp", port), MaxRetries: -1,
	}, nil)
	require.NoError(t, err)
	defer session.Close()

	expired, cancel := context.WithCancel(ctx)
	cancel()
	srv.Shutdown(expired) // the open SSE stream outlives the deadline; logged, not fatal
	<-done
}
