package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sns/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

var (
	mu        sync.Mutex
	published []string
)

// fakeSNS answers Publish like SNS (the query protocol) and keeps the
// messages published.
func fakeSNS(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil || r.Form.Get("Action") != "Publish" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mu.Lock()
	published = append(published, r.Form.Get("TopicArn")+" "+r.Form.Get("Message"))
	mu.Unlock()
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprint(w, `<PublishResponse xmlns="http://sns.amazonaws.com/doc/2010-03-31/"><PublishResult><MessageId>1</MessageId></PublishResult><ResponseMetadata><RequestId>1</RequestId></ResponseMetadata></PublishResponse>`)
}

// TestMain points the SNS client at fakeSNS. The custom endpoint goes
// in a config file: its key is the aws service id (SNS), and an env var
// would reach the factory in lower case.
func TestMain(m *testing.M) {
	ts := httptest.NewServer(http.HandlerFunc(fakeSNS))
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
        SNS:
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

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/messages", strings.NewReader(`{"text":"hi"}`))
	req.Header.Set("Content-Type", "application/json")
	srv.Instance().ServeHTTP(rec, req)
	require.Equal(t, http.StatusAccepted, rec.Code, rec.Body.String())

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/messages", strings.NewReader(`{"text":""}`))
	req.Header.Set("Content-Type", "application/json")
	srv.Instance().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code, "an empty message is refused")

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"arn:aws:sns:us-east-1:000000000000:example hi"}, published, "the message went to the topic example")
}
