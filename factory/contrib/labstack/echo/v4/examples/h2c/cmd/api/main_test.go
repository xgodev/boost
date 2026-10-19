package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/h2c/gen/builder/api"
	"go.uber.org/fx"
)

// TestMain reads the config of the example and lets echo listen on a free
// port: the test starts the app.
func TestMain(m *testing.M) {
	os.Setenv("CONF", "../../configs/config.yaml")
	os.Setenv("BOOST_FACTORY_ECHO_PORT", "0")
	boost.Start()
	os.Exit(m.Run())
}

func TestAPI(t *testing.T) {
	var rest *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&rest))...)
	require.NoError(t, app.Err())

	require.NoError(t, app.Start(context.Background()))
	t.Cleanup(func() { _ = app.Stop(context.Background()) })
	require.Eventually(t, func() bool {
		return rest.Instance().ListenerAddr() != nil
	}, 5*time.Second, 10*time.Millisecond, "the app serves echo")

	// HTTP/2 with prior knowledge, as curl --http2-prior-knowledge
	var protocols http.Protocols
	protocols.SetUnencryptedHTTP2(true)
	client := &http.Client{Transport: &http.Transport{Protocols: &protocols}}

	resp, err := client.Get("http://" + rest.Instance().ListenerAddr().String() + "/")
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "HTTP/2.0", resp.Proto)
	assert.Contains(t, string(body), "Protocol: HTTP/2.0")
	assert.Contains(t, string(body), "Method: GET")
}
