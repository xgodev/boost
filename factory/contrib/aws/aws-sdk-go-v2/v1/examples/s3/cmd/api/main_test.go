package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/s3/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

// fakeS3 answers HeadObject like S3: bucket example keeps only examplefile.
func fakeS3(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodHead || r.URL.Path != "/example/examplefile" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Last-Modified", "Mon, 02 Jan 2006 15:04:05 GMT")
	w.WriteHeader(http.StatusOK)
}

// TestMain points the S3 client at fakeS3. The custom endpoint goes
// in a config file: its key is the aws service id (S3), and an env var
// would reach the factory in lower case.
func TestMain(m *testing.M) {
	ts := httptest.NewServer(http.HandlerFunc(fakeS3))
	conf, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(conf, `boost:
  factory:
    aws:
      defaultRegion: us-east-1
      accessKeyId: test
      secretAccessKey: test
      customEndpoint:
        S3:
          url: %s
          hostnameImmutable: true
`, ts.URL)
	conf.Close()
	os.Setenv("CONF", conf.Name())
	boost.Start()
	code := m.Run()
	ts.Close()
	os.Remove(conf.Name())
	os.Exit(code)
}

func do(srv *boostecho.Server, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	rec := do(srv, "/files/examplefile")
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	assert.JSONEq(t, `{"name":"examplefile","lastModified":"2006-01-02T15:04:05Z"}`, rec.Body.String())

	rec = do(srv, "/files/missing")
	assert.Equal(t, http.StatusNotFound, rec.Code, rec.Body.String())
}
