package main

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/consumer/gen/builder/worker"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/consumer/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/consumer/internal/pkg/domain/service/message/receiver"
	"go.uber.org/fx"
)

var cluster *kafka.MockCluster

// TestMain points the consumer at a kafka cluster librdkafka runs in process
// and loads configs/config.yaml, which subscribes the topic messages.
func TestMain(m *testing.M) {
	var err error
	if cluster, err = kafka.NewMockCluster(1); err != nil {
		panic(err)
	}
	os.Setenv("BOOST_FACTORY_CONFLUENT_BROKERS", cluster.BootstrapServers())
	os.Setenv("CONF", "../../configs/config.yaml")
	boost.Start()
	os.Exit(m.Run())
}

// recorder keeps the messages the service of the app receives.
type recorder struct {
	receiver.Service
	mu       sync.Mutex
	received []string
}

func (r *recorder) Receive(ctx context.Context, msg *entity.Message) error {
	r.mu.Lock()
	r.received = append(r.received, msg.Text)
	r.mu.Unlock()
	return r.Service.Receive(ctx, msg)
}

func (r *recorder) got() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.received...)
}

// TestWorker starts app worker and produces a record on the topic: the kafka
// adapter hands it to the function, which hands it to the service.
func TestWorker(t *testing.T) {
	// the function runs the default command, kafka_confluent, when the
	// command line names none: drop the flags of go test.
	os.Args = os.Args[:1]

	rec := &recorder{}
	app := fx.New(append(worker.Modules(), fx.NopLogger,
		fx.Decorate(func(s receiver.Service) receiver.Service { rec.Service = s; return rec }),
	)...)
	require.NoError(t, app.Err())
	require.NoError(t, app.Start(context.Background()))

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cluster.BootstrapServers()})
	require.NoError(t, err)
	defer producer.Close()
	topic := "messages"
	require.NoError(t, producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("hi"),
	}, nil))
	producer.Flush(5000)

	require.Eventually(t, func() bool { return len(rec.got()) > 0 }, 20*time.Second, 100*time.Millisecond,
		"the function got no message from the topic")
	require.Equal(t, []string{"hi"}, rec.got())
}
