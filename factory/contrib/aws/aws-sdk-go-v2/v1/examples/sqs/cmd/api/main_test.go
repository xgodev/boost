package main

import (
	"crypto/md5"
	"encoding/hex"
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
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sqs/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

var (
	mu   sync.Mutex
	sent []string
)

// fakeSQS answers GetQueueUrl and SendMessage like SQS (the json protocol)
// and keeps the bodies sent.
func fakeSQS(w http.ResponseWriter, r *http.Request) {
	var in map[string]any
	_ = json.NewDecoder(r.Body).Decode(&in)
	w.Header().Set("Content-Type", "application/x-amz-json-1.0")
	switch r.Header.Get("X-Amz-Target") {
	case "AmazonSQS.GetQueueUrl":
		_ = json.NewEncoder(w).Encode(map[string]string{"QueueUrl": "http://" + r.Host + "/000000000000/" + in["QueueName"].(string)})
	case "AmazonSQS.SendMessage":
		body := in["MessageBody"].(string)
		mu.Lock()
		sent = append(sent, in["QueueUrl"].(string)+" "+body)
		mu.Unlock()
		sum := md5.Sum([]byte(body))
		_ = json.NewEncoder(w).Encode(map[string]string{"MessageId": "1", "MD5OfMessageBody": hex.EncodeToString(sum[:])})
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

// TestMain points the SQS client at fakeSQS. The custom endpoint goes
// in a config file: its key is the aws service id (SQS), and an env var
// would reach the factory in lower case.
func TestMain(m *testing.M) {
	ts := httptest.NewServer(http.HandlerFunc(fakeSQS))
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
        SQS:
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
	require.Len(t, sent, 1)
	assert.True(t, strings.HasSuffix(sent[0], "/000000000000/example hi"), "the message went to the queue example: %s", sent[0])
}
