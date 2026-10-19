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
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client/examples/greet/gen/builder/api"
	"go.uber.org/fx"
)

type greetIn struct {
	Name string `json:"name"`
}

type greetOut struct {
	Greeting string `json:"greeting"`
}

// TestMain points the MCP client at a greet server of the test, which
// answers as server/examples/greet.
func TestMain(m *testing.M) {
	s := mcp.NewServer(&mcp.Implementation{Name: "greet", Version: "1.0.0"}, nil)
	mcp.AddTool(s, &mcp.Tool{Name: "greet"}, func(_ context.Context, _ *mcp.CallToolRequest, in greetIn) (*mcp.CallToolResult, greetOut, error) {
		return nil, greetOut{Greeting: "hello " + in.Name}, nil
	})
	ts := httptest.NewServer(mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return s }, nil))
	os.Setenv("BOOST_FACTORY_MCP_CLIENT_ENDPOINT", ts.URL)
	boost.Start()
	code := m.Run()
	// The MCP session holds a stream open: drop it so Close does not wait.
	ts.CloseClientConnections()
	ts.Close()
	os.Exit(code)
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/greetings/boost", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"greeting":"hello boost"}`, rec.Body.String(), "the handler called the greet tool")
}
