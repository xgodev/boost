package main

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/pstest"
	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function/examples/pubsub/gen/builder/worker"
	"github.com/xgodev/boost/bootstrap/function/examples/pubsub/internal/app/worker/pubsub/handler/message/replier"
	"go.uber.org/fx"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const project = "example"

var srv *pstest.Server

// TestMain starts an in-memory pubsub server, points the pubsub factory at
// it and loads configs/config.yaml, which makes the adapter the default
// command.
func TestMain(m *testing.M) {
	srv = pstest.NewServer()
	os.Setenv("BOOST_FACTORY_GCP_PUBSUB_API__OPTIONS_EMULATOR__HOST", srv.Addr)
	os.Setenv("CONF", "../../configs/config.yaml")
	boost.Start()
	code := m.Run()
	srv.Close()
	os.Exit(code)
}

// topic creates topic id and a subscription sub to it.
func topic(t *testing.T, client *pubsub.Client, id, sub string) {
	ctx := context.Background()
	tp, err := client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{Name: "projects/" + project + "/topics/" + id})
	require.NoError(t, err)
	_, err = client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:  "projects/" + project + "/subscriptions/" + sub,
		Topic: tp.GetName(),
	})
	require.NoError(t, err)
}

// TestWorker starts app worker, publishes a message on topic messages and
// reads the reply the publisher middleware publishes on topic replies.
func TestWorker(t *testing.T) {
	// the function runs the default command, gcp_pubsub, when the command
	// line names none: drop the flags of go test.
	os.Args = os.Args[:1]

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project, option.WithEndpoint(srv.Addr), option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	require.NoError(t, err)
	defer client.Close()
	topic(t, client, "messages", "messages-worker")
	topic(t, client, replier.RepliesSubject, "replies-test")

	app := fx.New(append(worker.Modules(), fx.NopLogger)...)
	require.NoError(t, app.Err())
	require.NoError(t, app.Start(ctx))

	_, err = client.Publisher("messages").Publish(ctx, &pubsub.Message{Data: []byte(`{"message":"hi"}`)}).Get(ctx)
	require.NoError(t, err)

	rctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var reply *pubsub.Message
	err = client.Subscriber("replies-test").Receive(rctx, func(_ context.Context, msg *pubsub.Message) {
		msg.Ack()
		reply = msg
		cancel()
	})
	require.NoError(t, err)
	require.NotNil(t, reply, "the publisher middleware published no reply")
	assert.Equal(t, replier.RepliesSubject, reply.Attributes["ce_subject"])
	assert.JSONEq(t, `{"message":"hello, hi"}`, string(reply.Data))
}
