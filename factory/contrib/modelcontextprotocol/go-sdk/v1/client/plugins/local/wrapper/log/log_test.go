package log

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	"github.com/xgodev/boost/wrapper/log"
)

type recorder struct {
	log.Noop
	mu    *sync.Mutex
	lines *[]string
}

func (r recorder) FromContext(context.Context) log.Logger { return r }
func (r recorder) WithField(string, any) log.Logger       { return r }
func (r recorder) Infof(f string, _ ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	*r.lines = append(*r.lines, f)
}

func (r recorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(*r.lines)
}

// connect loads config from env (boost's koanf keeps keys across Load calls,
// so each test sets the keys it reads) and connects a client with plugin.
func connect(t *testing.T, env map[string]string, plugin client.Plugin) (*client.Client, recorder, error) {
	t.Helper()
	for k, v := range env {
		t.Setenv(k, v)
	}
	config.Set(koanf.New())
	config.Load()

	rec := recorder{mu: new(sync.Mutex), lines: new([]string)}
	log.Set(rec)
	t.Cleanup(func() { log.Set(log.Noop{}) })

	srv := mcp.NewServer(&mcp.Implementation{Name: "upstream", Version: "1.0.0"}, nil)
	ts := httptest.NewServer(mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return srv }, nil))
	t.Cleanup(ts.Close)

	cli, err := client.NewClientWithOptions(context.Background(), &client.Options{
		Name: "test", Version: "1.0.0", Transport: client.TransportHTTP, Endpoint: ts.URL, MaxRetries: -1,
	}, plugin)
	if err == nil {
		t.Cleanup(func() { cli.Close() })
	}
	return cli, rec, err
}

func TestRegister_LogsCallsAtConfiguredLevel(t *testing.T) {
	cli, rec, err := connect(t, map[string]string{
		"BOOST_FACTORY_MCP_CLIENT_PLUGINS_LOG_ENABLED": "true",
		"BOOST_FACTORY_MCP_CLIENT_PLUGINS_LOG_LEVEL":   "INFO",
	}, Register)
	require.NoError(t, err)
	before := rec.count()

	_, err = cli.Session().ListTools(context.Background(), nil)
	require.NoError(t, err)

	assert.Positive(t, before, "initialize is logged: the plugin runs before Connect")
	assert.Greater(t, rec.count(), before)
}

func TestRegister_Disabled(t *testing.T) {
	cli, rec, err := connect(t, map[string]string{"BOOST_FACTORY_MCP_CLIENT_PLUGINS_LOG_ENABLED": "false"}, Register)
	require.NoError(t, err)

	_, err = cli.Session().ListTools(context.Background(), nil)
	require.NoError(t, err)
	assert.Zero(t, rec.count())
}

func TestRegister_OptionsError(t *testing.T) {
	_, _, err := connect(t, map[string]string{"BOOST_FACTORY_MCP_CLIENT_PLUGINS_LOG_ENABLED": "maybe"}, Register)
	assert.Error(t, err)
}

func TestNewOptionsWithPath(t *testing.T) {
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_PLUGINS_LOG_ENABLED", "true")
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_GITHUB_PLUGINS_LOG_LEVEL", "TRACE")
	config.Set(koanf.New())
	config.Load()

	o, err := NewOptionsWithPath("boost.factory.mcp.client.github.plugins.log")
	require.NoError(t, err)
	assert.True(t, o.Enabled)
	assert.Equal(t, "TRACE", o.Level)
}
