package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/producer/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

var cluster *kafka.MockCluster

// TestMain points the producer at a kafka cluster librdkafka runs in process.
func TestMain(m *testing.M) {
	var err error
	if cluster, err = kafka.NewMockCluster(1); err != nil {
		panic(err)
	}
	os.Setenv("BOOST_FACTORY_CONFLUENT_BROKERS", cluster.BootstrapServers())
	boost.Start()
	code := m.Run()
	cluster.Close()
	os.Exit(code)
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

	rec := post(srv, `{"text":"hi"}`)
	require.Equal(t, http.StatusAccepted, rec.Code, rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, post(srv, `{"text":""}`).Code, "an empty message is refused")

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cluster.BootstrapServers(),
		"group.id":          "test",
		"auto.offset.reset": "earliest",
	})
	require.NoError(t, err)
	defer consumer.Close()
	require.NoError(t, consumer.Subscribe("messages", nil))
	msg, err := consumer.ReadMessage(10 * time.Second)
	require.NoError(t, err)
	assert.Equal(t, "hi", string(msg.Value), "the message went to the topic messages")
}
