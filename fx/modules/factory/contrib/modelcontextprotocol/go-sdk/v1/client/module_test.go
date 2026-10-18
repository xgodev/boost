package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModule(t *testing.T) {
	upstream := mcp.NewServer(&mcp.Implementation{Name: "upstream", Version: "1.0.0"}, nil)
	ts := httptest.NewServer(mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return upstream }, nil))
	defer ts.Close()
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_ENDPOINT", ts.URL)
	boost.Start()

	var (
		cli     *client.Client
		session *mcp.ClientSession
	)
	app := fxtest.New(t, Module(), fx.Populate(&cli, &session)).RequireStart()

	assert.Same(t, cli.Session(), session)
	_, err := session.ListTools(context.Background(), nil)
	require.NoError(t, err)

	app.RequireStop()
	_, err = session.ListTools(context.Background(), nil)
	assert.Error(t, err, "the session is closed on app stop")
}
