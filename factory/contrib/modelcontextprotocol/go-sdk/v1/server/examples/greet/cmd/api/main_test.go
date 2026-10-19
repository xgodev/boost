package main

import (
	"context"
	"os"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	mcpserver "github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server/examples/greet/gen/builder/api"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func TestTools(t *testing.T) {
	var srv *mcpserver.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := srv.Instance().Connect(ctx, serverTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = serverSession.Close() })

	session, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil).Connect(ctx, clientTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })

	tools, err := session.ListTools(ctx, nil)
	require.NoError(t, err)
	require.Len(t, tools.Tools, 1)
	assert.Equal(t, "greet", tools.Tools[0].Name)
	assert.Equal(t, "Says hello.", tools.Tools[0].Description)

	res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "greet", Arguments: map[string]any{"name": "boost"}})
	require.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Equal(t, map[string]any{"greeting": "hello boost"}, res.StructuredContent)
}
