package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/examples/cobra/gen/builder/api"
	mcpserver "github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	"go.uber.org/fx"
)

var google, bing *httptest.Server

// TestMain reads the config of the example and points each client at a
// server of the test that answers with its own status.
func TestMain(m *testing.M) {
	google = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	bing = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	os.Setenv("CONF", "../../configs/config.yaml")
	os.Setenv("BOOST_FACTORY_RESTY_CLIENTS_GOOGLE_HOST", google.URL)
	os.Setenv("BOOST_FACTORY_RESTY_CLIENTS_BING_HOST", bing.URL)
	boost.Start()
	code := m.Run()
	google.Close()
	bing.Close()
	os.Exit(code)
}

func TestTools(t *testing.T) {
	var srv *mcpserver.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())
	assert.Equal(t, mcpserver.TransportSTDIO, srv.Options().Transport, "the config runs the server on stdio")

	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := srv.Instance().Connect(ctx, serverTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = serverSession.Close() })

	session, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil).Connect(ctx, clientTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })

	for tool, want := range map[string]map[string]any{
		"google": {"url": google.URL, "status": float64(http.StatusOK)},
		"bing":   {"url": bing.URL, "status": float64(http.StatusAccepted)},
	} {
		res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tool, Arguments: map[string]any{}})
		require.NoError(t, err)
		assert.False(t, res.IsError, tool)
		assert.Equal(t, want, res.StructuredContent, "each tool calls its own client")
	}
}
