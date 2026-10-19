package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function/examples/middlewares/gen/builder/worker"
	"go.uber.org/fx"
)

var url string

// TestMain puts the adapter on a free port and loads configs/config.yaml,
// which makes it the default command.
func TestMain(m *testing.M) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	url = fmt.Sprintf("http://127.0.0.1:%d/", port)
	os.Setenv("BOOST_BOOTSTRAP_FUNCTION_ADAPTER_CLOUDEVENTS_HTTP_PORT", strconv.Itoa(port))
	os.Setenv("CONF", "../../configs/config.yaml")
	boost.Start()
	os.Exit(m.Run())
}

// post posts an event with data and returns the status of the answer.
func post(data string) (int, error) {
	body := fmt.Sprintf(`{"specversion":"1.0","id":"1","source":"test","type":"number.sent","datacontenttype":"application/json","data":%s}`, data)
	res, err := http.Post(url, "application/cloudevents+json", strings.NewReader(body))
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	return res.StatusCode, nil
}

// TestWorker starts app worker and posts events to the adapter: the event
// without number panics, and the worker answers 500 and keeps echoing.
func TestWorker(t *testing.T) {
	// the function runs the default command, cloudevents_http, when the
	// command line names none: drop the flags of go test.
	os.Args = os.Args[:1]

	app := fx.New(append(worker.Modules(), fx.NopLogger)...)
	require.NoError(t, app.Err())
	require.NoError(t, app.Start(context.Background()))

	var status int
	require.Eventually(t, func() bool {
		var err error
		status, err = post(`{"number":2}`)
		return err == nil
	}, 10*time.Second, 50*time.Millisecond, "the adapter did not listen")
	assert.Equal(t, http.StatusOK, status)

	status, err := post(`{}`)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, status, "recovery turns the panic into an error")

	status, err = post(`{"number":3}`)
	require.NoError(t, err, "the worker survives the panic")
	assert.Equal(t, http.StatusOK, status)

	assert.Equal(t, 2.0, processed(t, "success"), "prometheus counts the echoes")
	assert.Equal(t, 1.0, processed(t, "error"), "prometheus counts the recovered panic")
}

// processed is the count of boost_function_messages_processed_total with
// label status.
func processed(t *testing.T, status string) float64 {
	families, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)
	for _, family := range families {
		if family.GetName() != "boost_function_messages_processed_total" {
			continue
		}
		for _, metric := range family.GetMetric() {
			for _, label := range metric.GetLabel() {
				if label.GetName() == "status" && label.GetValue() == status {
					return metric.GetCounter().GetValue()
				}
			}
		}
	}
	return 0
}
