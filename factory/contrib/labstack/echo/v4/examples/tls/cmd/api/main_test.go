package main

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/tls/gen/builder/api"
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
		return rest.Instance().TLSListenerAddr() != nil
	}, 5*time.Second, 10*time.Millisecond, "the app serves echo over TLS")

	// the certificate is generated at the start, as curl -k
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // self-signed certificate of the example
		ForceAttemptHTTP2: true,
	}}

	resp, err := client.Get("https://" + rest.Instance().TLSListenerAddr().String() + "/")
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.TLS)
	assert.Contains(t, string(body), "Method: GET")
	assert.Contains(t, string(body), "Path: /")
}
