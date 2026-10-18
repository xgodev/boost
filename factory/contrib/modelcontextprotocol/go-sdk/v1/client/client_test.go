package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
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

// stdioServerEnv re-runs this test binary as an MCP server on stdio, which is
// what the COMMAND transport spawns in TestNewClientWithOptions_Command.
const stdioServerEnv = "BOOST_MCP_CLIENT_TEST_STDIO_SERVER"

func TestMain(m *testing.M) {
	if os.Getenv(stdioServerEnv) == "1" {
		if err := upstreamServer().Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
	os.Exit(m.Run())
}

type echoIn struct {
	Text string `json:"text"`
}

type echoOut struct {
	Text string `json:"text"`
}

func upstreamServer() *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{Name: "upstream", Version: "1.0.0"}, nil)
	mcp.AddTool(srv, &mcp.Tool{Name: "echo"},
		func(_ context.Context, _ *mcp.CallToolRequest, in echoIn) (*mcp.CallToolResult, echoOut, error) {
			return nil, echoOut{Text: in.Text}, nil
		})
	return srv
}

func httpUpstream(t *testing.T) string {
	t.Helper()
	srv := upstreamServer()
	ts := httptest.NewServer(mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return srv }, nil))
	t.Cleanup(ts.Close)
	return ts.URL
}

func loadConfig(t *testing.T) {
	t.Helper()
	config.Set(koanf.New())
	config.Load()
}

func callEcho(t *testing.T, cli *Client) {
	t.Helper()
	res, err := cli.Session().CallTool(context.Background(), &mcp.CallToolParams{
		Name: "echo", Arguments: map[string]any{"text": "hi"},
	})
	require.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Equal(t, map[string]any{"text": "hi"}, res.StructuredContent)
}

func httpOptions(endpoint string) *Options {
	return &Options{Name: "test", Version: "1.0.0", Transport: TransportHTTP, Endpoint: endpoint, MaxRetries: -1}
}

func TestNewOptions_Defaults(t *testing.T) {
	loadConfig(t)

	o, err := NewOptions()
	require.NoError(t, err)
	assert.Equal(t, boost.ApplicationName(), o.Name)
	assert.Equal(t, "1.0.0", o.Version)
	assert.Equal(t, TransportHTTP, o.Transport)
	assert.Equal(t, "http://localhost:8090/mcp", o.Endpoint)
	assert.Equal(t, 5, o.MaxRetries)
	assert.Empty(t, o.Command)
	assert.Empty(t, o.Args)
	assert.Zero(t, o.KeepAlive)
}

func TestNewOptionsWithPath_Overrides(t *testing.T) {
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_GITHUB_TRANSPORT", "COMMAND")
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_GITHUB_COMMAND", "github-mcp-server")
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_GITHUB_KEEP__ALIVE", "15s")
	ConfigAdd(root + ".github")
	loadConfig(t)

	o, err := NewOptionsWithPath(root + ".github")
	require.NoError(t, err)
	assert.Equal(t, TransportCOMMAND, o.Transport)
	assert.Equal(t, "github-mcp-server", o.Command)
	assert.Equal(t, 15*time.Second, o.KeepAlive)
	assert.Equal(t, 5, o.MaxRetries)
}

func TestNewClient_FromConfig(t *testing.T) {
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_ENDPOINT", httpUpstream(t))
	loadConfig(t)
	ctx := context.Background()

	cli, err := NewClient(ctx)
	require.NoError(t, err)
	defer cli.Close()
	callEcho(t, cli)

	cli2, err := NewClientWithConfigPath(ctx, root)
	require.NoError(t, err)
	defer cli2.Close()
	assert.NotNil(t, cli2.Instance())
	assert.Equal(t, TransportHTTP, cli2.Options().Transport)
}

func TestNewClient_ConfigErrorIsReturned(t *testing.T) {
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_MAX__RETRIES", "many")
	loadConfig(t)
	ctx := context.Background()

	_, err := NewClient(ctx)
	assert.Error(t, err)
	_, err = NewClientWithConfigPath(ctx, root)
	assert.Error(t, err)
}

func TestNewClientWithOptions_HTTP(t *testing.T) {
	cli, err := NewClientWithOptions(context.Background(), httpOptions(httpUpstream(t)))
	require.NoError(t, err)
	callEcho(t, cli)
	assert.NoError(t, cli.Close())
}

func TestNewClientWithOptions_Command(t *testing.T) {
	t.Setenv(stdioServerEnv, "1")
	o := &Options{Name: "test", Version: "1.0.0", Transport: "command", Command: os.Args[0], Args: []string{"-test.run=^$"}}

	cli, err := NewClientWithOptions(context.Background(), o)
	require.NoError(t, err)
	callEcho(t, cli)
	assert.NoError(t, cli.Close())
}

func TestNewClientWithOptions_CommandRequired(t *testing.T) {
	o := &Options{Name: "test", Version: "1.0.0", Transport: TransportCOMMAND}

	_, err := NewClientWithOptions(context.Background(), o)
	assert.True(t, errors.IsNotValid(err), "got %v", err)
}

func TestNewClientWithOptions_UnknownTransport(t *testing.T) {
	o := httpOptions("http://localhost")
	o.Transport = "SSE"

	_, err := NewClientWithOptions(context.Background(), o)
	assert.True(t, errors.IsNotValid(err), "got %v", err)
}

func TestNewClientWithOptions_ConnectError(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	endpoint := fmt.Sprintf("http://%s/mcp", l.Addr())
	l.Close()

	_, err = NewClientWithOptions(context.Background(), httpOptions(endpoint))
	assert.Error(t, err)
}

func TestNewClientWithOptions_PluginsRunBeforeConnect(t *testing.T) {
	ctx := context.Background()
	var connected bool
	plugin := func(_ context.Context, c *Client) error {
		connected = c.Session() != nil
		c.Instance().AddSendingMiddleware()
		return nil
	}

	cli, err := NewClientWithOptions(ctx, httpOptions(httpUpstream(t)), plugin)
	require.NoError(t, err)
	defer cli.Close()
	assert.False(t, connected, "plugins must run before Connect")

	boom := fmt.Errorf("boom")
	_, err = NewClientWithOptions(ctx, httpOptions(httpUpstream(t)), func(context.Context, *Client) error { return boom })
	assert.ErrorIs(t, err, boom)
}
