package nats

import (
	"context"
	"fmt"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1"
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"

	"github.com/stretchr/testify/assert"
)

const TestPort = 8369

type Handler struct {
}

func (h *Handler) Handle(ctx context.Context, in v2.Event) (*v2.Event, error) {
	return nil, nil
}

func runServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	return runServerWithOptions(&opts)
}

func runServerWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}

func TestSubscriberListenerSubscribe(t *testing.T) {

	boost.Start()
	var err error
	var options *nats.Options

	s := runServerOnPort(TestPort)
	defer s.Shutdown()

	sUrl := fmt.Sprintf("nats://127.0.0.1:%d", TestPort)

	options, err = nats.NewOptions()
	assert.Nil(t, err)

	options.Url = sUrl

	conn, err := nats.NewConnWithOptions(context.Background(), options)
	assert.Nil(t, err)

	lis := NewSubscriber[*v2.Event](conn, nil, "subject", "queue")
	subscribe, err := lis.Subscribe(context.Background())
	assert.Nil(t, err)

	assert.True(t, subscribe.IsValid())

	err = subscribe.Unsubscribe()
	assert.Nil(t, err)
}
