package swagger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggo/swag"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const spec = `{"swagger":"2.0","info":{"title":"boost","version":"1.0"}}`

type doc struct{}

func (doc) ReadDoc() string { return spec }

func TestMain(m *testing.M) {
	boost.Start()
	swag.Register(swag.Name, doc{})
	os.Exit(m.Run())
}

func newServer(t *testing.T) *echo.Server {
	t.Helper()
	srv, err := echo.NewServerWithOptions(context.Background(), &echo.Options{Type: "REST"})
	require.NoError(t, err)
	return srv
}

func get(srv *echo.Server, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestRegisterServesTheSwaggerUI(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, Register(context.Background(), srv))

	rec := get(srv, "/swagger/index.html")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "swagger-ui")
}

func TestRegisterServesTheRegisteredDoc(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, Register(context.Background(), srv))

	rec := get(srv, "/swagger/doc.json")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, spec, rec.Body.String())
}

func TestRegisterUsesTheConfiguredRoute(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, NewSwaggerWithOptions(&Options{Enabled: true, Route: "/docs/*"}).Register(context.Background(), srv))

	assert.Equal(t, http.StatusOK, get(srv, "/docs/index.html").Code)
	assert.Equal(t, http.StatusNotFound, get(srv, "/swagger/index.html").Code)
}

func TestRegisterDisabledAddsNoRoute(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, NewSwaggerWithOptions(&Options{Route: "/swagger/*"}).Register(context.Background(), srv))

	assert.Equal(t, http.StatusNotFound, get(srv, "/swagger/index.html").Code)
}

func TestNewSwaggerReadsTheDefaults(t *testing.T) {
	s, err := NewSwagger()

	require.NoError(t, err)
	assert.Equal(t, &Options{Enabled: true, Route: "/swagger/*"}, s.options)
}

func TestNewSwaggerWithConfigPathReadsThePath(t *testing.T) {
	s, err := NewSwaggerWithConfigPath(root)

	require.NoError(t, err)
	assert.Equal(t, &Options{Enabled: true, Route: "/swagger/*"}, s.options)
}

func TestConstructorsReturnTheConfigError(t *testing.T) {
	const env = "BOOST_FACTORY_ECHO_PLUGINS_SWAGGER_ENABLED"
	t.Setenv(env, "not-a-bool")
	config.Load()
	// koanf merges every load and keeps a key its env var no longer sets, so
	// the default goes back by loading it before t.Setenv unsets the var.
	t.Cleanup(func() {
		if err := os.Setenv(env, "true"); err != nil {
			t.Error(err)
		}
		config.Load()
	})

	_, err := NewSwagger()
	assert.Error(t, err)

	_, err = NewSwaggerWithConfigPath(root)
	assert.Error(t, err)

	assert.Error(t, Register(context.Background(), newServer(t)))
}
