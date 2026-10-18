package log

import (
	"context"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
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

// setup loads config from env. boost's koanf keeps keys across Load calls, so
// every test sets the keys it depends on instead of trusting the defaults.
func setup(t *testing.T, env map[string]string) (*server.Server, recorder) {
	t.Helper()
	for k, v := range env {
		t.Setenv(k, v)
	}
	config.Set(koanf.New())
	config.Load()

	rec := recorder{mu: new(sync.Mutex), lines: new([]string)}
	log.Set(rec)
	t.Cleanup(func() { log.Set(log.Noop{}) })

	srv, err := server.NewServerWithOptions(context.Background(),
		&server.Options{Name: "test", Version: "1.0.0", Transport: server.TransportHTTP})
	require.NoError(t, err)
	return srv, rec
}

func listTools(t *testing.T, srv *server.Server) {
	t.Helper()
	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()
	ss, err := srv.Instance().Connect(ctx, st, nil)
	require.NoError(t, err)
	defer ss.Close()

	cs, err := mcp.NewClient(&mcp.Implementation{Name: "c", Version: "1.0.0"}, nil).Connect(ctx, ct, nil)
	require.NoError(t, err)
	defer cs.Close()
	_, err = cs.ListTools(ctx, nil)
	require.NoError(t, err)
}

func TestRegister_LogsCallsAtConfiguredLevel(t *testing.T) {
	srv, rec := setup(t, map[string]string{
		"BOOST_FACTORY_MCP_SERVER_PLUGINS_LOG_ENABLED": "true",
		"BOOST_FACTORY_MCP_SERVER_PLUGINS_LOG_LEVEL":   "INFO",
	})

	require.NoError(t, Register(context.Background(), srv))
	listTools(t, srv)

	assert.Positive(t, rec.count())
}

func TestRegister_Disabled(t *testing.T) {
	srv, rec := setup(t, map[string]string{"BOOST_FACTORY_MCP_SERVER_PLUGINS_LOG_ENABLED": "false"})

	require.NoError(t, Register(context.Background(), srv))
	listTools(t, srv)

	assert.Zero(t, rec.count())
}

func TestRegister_OptionsError(t *testing.T) {
	srv, _ := setup(t, map[string]string{"BOOST_FACTORY_MCP_SERVER_PLUGINS_LOG_ENABLED": "maybe"})

	assert.Error(t, Register(context.Background(), srv))
}

func TestNewOptionsWithPath(t *testing.T) {
	t.Setenv("BOOST_FACTORY_MCP_SERVER_PLUGINS_LOG_ENABLED", "true")
	t.Setenv("BOOST_FACTORY_MCP_SERVER_ADMIN_PLUGINS_LOG_LEVEL", "TRACE")
	config.Set(koanf.New())
	config.Load()

	o, err := NewOptionsWithPath("boost.factory.mcp.server.admin.plugins.log")
	require.NoError(t, err)
	assert.True(t, o.Enabled)
	assert.Equal(t, "TRACE", o.Level)
}
