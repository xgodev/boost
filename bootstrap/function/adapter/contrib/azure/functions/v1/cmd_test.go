package azure

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func TestNewIsTheAzureCommand(t *testing.T) {
	assert.Equal(t, "azure", New[*event.Event]()(nil).Use)
}

func TestNewOptionsReadsTheConfig(t *testing.T) {
	o, err := NewOptions()
	require.NoError(t, err)
	assert.Equal(t, "$return", o.Output)
	assert.Equal(t, "api", o.RoutePrefix)
	assert.Equal(t, 10, o.MaxWorkers)
	assert.Positive(t, o.Port)
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func TestServeAnswersUntilTheContextEnds(t *testing.T) {
	o := options()
	o.Port = freePort(t)
	h := NewHandler[*event.Event](func(context.Context, cloudevents.Event) (*event.Event, error) { return nil, nil }, o)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Serve(ctx, h, o.Port) }()

	url := "http://127.0.0.1:" + strconv.Itoa(o.Port) + "/"
	require.Eventually(t, func() bool {
		res, err := http.Get(url)
		if err != nil {
			return false
		}
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 2*time.Second, 10*time.Millisecond)

	cancel()
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Serve did not return after the context ended")
	}
}
