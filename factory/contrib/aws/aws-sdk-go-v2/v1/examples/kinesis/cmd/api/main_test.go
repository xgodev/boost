package main

import (
	"encoding/base64"
	"encoding/json"
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
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/kinesis/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

var (
	mu  sync.Mutex
	put []string
)

// fakeKinesis answers PutRecord like Kinesis (the json protocol) and keeps
// the records put.
func fakeKinesis(w http.ResponseWriter, r *http.Request) {
	var in struct{ StreamName, PartitionKey, Data string }
	if r.Header.Get("X-Amz-Target") != "Kinesis_20131202.PutRecord" || json.NewDecoder(r.Body).Decode(&in) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := base64.StdEncoding.DecodeString(in.Data)
	mu.Lock()
	put = append(put, in.StreamName+" "+in.PartitionKey+" "+string(data))
	mu.Unlock()
	w.Header().Set("Content-Type", "application/x-amz-json-1.1")
	fmt.Fprint(w, `{"SequenceNumber":"1","ShardId":"shardId-000000000000"}`)
}

// TestMain points the Kinesis client at fakeKinesis. The custom endpoint goes
// in a config file: its key is the aws service id (Kinesis), and an env var
// would reach the factory in lower case.
func TestMain(m *testing.M) {
	ts := httptest.NewServer(http.HandlerFunc(fakeKinesis))
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
        Kinesis:
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
	assert.Equal(t, []string{"example messages hi"}, put, "the message went to the stream example")
}
