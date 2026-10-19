package kafka

import (
	"context"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	factory "github.com/xgodev/boost/factory/contrib/segmentio/kafka-go/v0"
)

func TestNew_Command(t *testing.T) {
	assert.Equal(t, "kafkago", New[*event.Event]()(nil).Use)
}

func TestNewOptions_FromEnv(t *testing.T) {
	key := "BOOST_BOOTSTRAP_FUNCTION_ADAPTER_KAFKAGO_TOPICS"
	t.Setenv(key, "orders")
	boost.Start()
	t.Cleanup(func() {
		os.Setenv(key, "changeme")
		boost.Start()
	})

	o, err := NewOptions()
	require.NoError(t, err)
	assert.Equal(t, []string{"orders"}, o.Topics)
	assert.Equal(t, -1, o.RetryLimit)
}

func TestNewReader_ReadsTheTopicsInTheGroup(t *testing.T) {
	fo := &factory.Options{
		Brokers: []string{"localhost:9092"},
		Reader:  factory.ReaderOptions{GroupID: "g", Topic: "ignored", StartOffset: "last"},
	}

	r, err := newReader(context.Background(), fo, []string{"a", "b"})
	require.NoError(t, err)
	t.Cleanup(func() { _ = r.Close() })
	assert.Equal(t, []string{"a", "b"}, r.Config().GroupTopics)
	assert.Empty(t, r.Config().Topic)
	assert.Equal(t, "g", r.Config().GroupID)
}

func TestNewReader_NeedsAGroup(t *testing.T) {
	fo := &factory.Options{Brokers: []string{"localhost:9092"}, Reader: factory.ReaderOptions{StartOffset: "last"}}

	_, err := newReader(context.Background(), fo, []string{"a"})
	assert.ErrorContains(t, err, "groupId")
}
