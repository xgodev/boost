package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/helloworld/gen/builder/api"
	"go.uber.org/fx"
)

const (
	childEnv = "EXAMPLE_CHILD"
	// messageEnv sets internal.pkg.domain.service.greeting.greeter.local.message.
	messageEnv = "INTERNAL_PKG_DOMAIN_SERVICE_GREETING_GREETER_LOCAL_MESSAGE"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// rest builds app api from its generated modules and returns its echo
// server, which serves requests without listening on a port.
func rest(t *testing.T) *boostecho.Server {
	t.Helper()
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())
	return srv
}

func get(srv *boostecho.Server, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestAPI(t *testing.T) {
	if os.Getenv(childEnv) == "1" {
		t.Skip("the child runs TestAPIMessageChild only")
	}
	srv := rest(t)

	hello := get(srv, "/hello-world")
	assert.Equal(t, http.StatusOK, hello.Code)
	assert.JSONEq(t, `{"message":"Hello World!!"}`, hello.Body.String())
	assert.NotEmpty(t, hello.Header().Get("X-Request-Id"), "the requestid plugin runs")

	assert.Equal(t, http.StatusOK, get(srv, "/health").Code, "the health plugin serves its route")
	assert.Equal(t, http.StatusOK, get(srv, "/resource-status").Code, "the restresponse plugin serves its route")
}

// TestAPIMessage boots app api in a new process with the message set from
// the environment: config is read once per process.
func TestAPIMessage(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestAPIMessageChild$", "-test.v")
	cmd.Env = append(os.Environ(), messageEnv+"=Hi", childEnv+"=1")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", out)
	assert.Contains(t, string(out), "--- PASS: TestAPIMessageChild")
}

func TestAPIMessageChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestAPIMessage")
	}
	assert.JSONEq(t, `{"message":"Hi"}`, get(rest(t), "/hello-world").Body.String())
}
