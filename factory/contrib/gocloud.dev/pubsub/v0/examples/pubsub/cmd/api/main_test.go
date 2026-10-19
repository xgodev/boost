package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0/examples/pubsub/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
	"gocloud.dev/pubsub"
)

// TestMain boots boost: the factory opens its default topic, mem://topicA.
func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func post(srv *boostecho.Server, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	srv.Instance().ServeHTTP(rec, req)
	return rec
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sub, err := pubsub.OpenSubscription(ctx, "mem://topicA")
	require.NoError(t, err)
	defer sub.Shutdown(ctx)

	rec := post(srv, `{"text":"hi"}`)
	require.Equal(t, http.StatusAccepted, rec.Code, rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, post(srv, `{"text":""}`).Code, "an empty message is refused")

	msg, err := sub.Receive(ctx)
	require.NoError(t, err)
	msg.Ack()
	assert.Equal(t, "hi", string(msg.Body), "the message went to the topic")
}
