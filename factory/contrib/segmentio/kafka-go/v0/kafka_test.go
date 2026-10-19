package kafka

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
)

func TestNewOptions_Defaults(t *testing.T) {
	boost.Start()

	o, err := NewOptions()
	require.NoError(t, err)

	assert.Equal(t, []string{"localhost:9092"}, o.Brokers)
	assert.Equal(t, "crc32", o.Writer.Balancer)
	assert.Equal(t, "all", o.Writer.RequiredAcks)
	assert.Equal(t, 10, o.Writer.MaxAttempts)
	assert.Equal(t, 10*time.Millisecond, o.Writer.BatchTimeout)
	assert.Equal(t, 100, o.Reader.QueueCapacity)
	assert.Equal(t, 10485760, o.Reader.MaxBytes)
	assert.Equal(t, "last", o.Reader.StartOffset)
	assert.Equal(t, "leader", o.Conn.Type)
}

func TestNewOptions_FromEnv(t *testing.T) {
	key := "BOOST_FACTORY_KAFKAGO_WRITER_REQUIRED__ACKS"
	t.Setenv(key, "one")
	boost.Start()
	// boost keeps config keys across loads: reload the default.
	t.Cleanup(func() {
		os.Setenv(key, "all")
		boost.Start()
	})

	o, err := NewOptions()
	require.NoError(t, err)
	assert.Equal(t, "one", o.Writer.RequiredAcks)
}

func writerOptions() *Options {
	return &Options{
		Brokers: []string{"b1:9092", "b2:9092"},
		Writer: WriterOptions{
			Balancer:     "hash",
			RequiredAcks: "one",
			Compression:  "snappy",
			MaxAttempts:  3,
			BatchSize:    50,
			BatchBytes:   2048,
			BatchTimeout: 5 * time.Millisecond,
			ReadTimeout:  time.Second,
			WriteTimeout: 2 * time.Second,
			Async:        true,
		},
	}
}

func TestNewWriterWithOptions_MapsOptions(t *testing.T) {
	w, err := NewWriterWithOptions(context.Background(), writerOptions())
	require.NoError(t, err)
	t.Cleanup(func() { _ = w.Close() })

	assert.Equal(t, "b1:9092,b2:9092", w.Addr.String())
	assert.Empty(t, w.Topic, "the topic comes from each message")
	assert.IsType(t, &kafka.Hash{}, w.Balancer)
	assert.Equal(t, kafka.RequireOne, w.RequiredAcks)
	assert.Equal(t, kafka.Snappy, w.Compression)
	assert.Equal(t, 3, w.MaxAttempts)
	assert.Equal(t, 50, w.BatchSize)
	assert.Equal(t, int64(2048), w.BatchBytes)
	assert.Equal(t, 5*time.Millisecond, w.BatchTimeout)
	assert.Equal(t, time.Second, w.ReadTimeout)
	assert.Equal(t, 2*time.Second, w.WriteTimeout)
	assert.True(t, w.Async)
	assert.NotNil(t, w.Logger)
	assert.NotNil(t, w.ErrorLogger)
}

func TestNewWriterWithOptions_RejectsUnknownValues(t *testing.T) {
	cases := map[string]func(o *Options){
		"balancer":     func(o *Options) { o.Writer.Balancer = "sticky" },
		"requiredAcks": func(o *Options) { o.Writer.RequiredAcks = "two" },
		"compression":  func(o *Options) { o.Writer.Compression = "brotli" },
		"no brokers":   func(o *Options) { o.Brokers = nil },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			o := writerOptions()
			change(o)
			_, err := NewWriterWithOptions(context.Background(), o)
			assert.Error(t, err)
		})
	}
}

func TestNewWriterWithOptions_AppliesPlugins(t *testing.T) {
	plugin := func(_ context.Context, w *kafka.Writer) error {
		w.AllowAutoTopicCreation = true
		return nil
	}
	w, err := NewWriterWithOptions(context.Background(), writerOptions(), plugin)
	require.NoError(t, err)
	t.Cleanup(func() { _ = w.Close() })
	assert.True(t, w.AllowAutoTopicCreation)

	failing := func(context.Context, *kafka.Writer) error { return errors.New("boom") }
	_, err = NewWriterWithOptions(context.Background(), writerOptions(), failing)
	assert.ErrorContains(t, err, "boom")
}

func readerOptions() *Options {
	return &Options{
		Brokers: []string{"b1:9092"},
		Reader: ReaderOptions{
			GroupID:          "group",
			Topic:            "topic",
			QueueCapacity:    10,
			MinBytes:         1,
			MaxBytes:         1024,
			MaxWait:          time.Second,
			ReadBatchTimeout: 2 * time.Second,
			CommitInterval:   3 * time.Second,
			StartOffset:      "first",
		},
	}
}

func TestNewReaderWithOptions_MapsOptions(t *testing.T) {
	r, err := NewReaderWithOptions(context.Background(), readerOptions())
	require.NoError(t, err)
	t.Cleanup(func() { _ = r.Close() })

	c := r.Config()
	assert.Equal(t, []string{"b1:9092"}, c.Brokers)
	assert.Equal(t, "group", c.GroupID)
	assert.Equal(t, "topic", c.Topic)
	assert.Equal(t, 10, c.QueueCapacity)
	assert.Equal(t, 1, c.MinBytes)
	assert.Equal(t, 1024, c.MaxBytes)
	assert.Equal(t, time.Second, c.MaxWait)
	assert.Equal(t, 2*time.Second, c.ReadBatchTimeout)
	assert.Equal(t, 3*time.Second, c.CommitInterval)
	assert.Equal(t, kafka.FirstOffset, c.StartOffset)
	assert.NotNil(t, c.Logger)
	assert.NotNil(t, c.ErrorLogger)
}

func TestNewReaderWithOptions_InvalidConfigIsAnError(t *testing.T) {
	cases := map[string]func(o *Options){
		// kafka.NewReader panics on these; the factory must return an error.
		"no topic":            func(o *Options) { o.Reader.Topic = "" },
		"group and partition": func(o *Options) { o.Reader.Partition = 1 },
		"no brokers":          func(o *Options) { o.Brokers = nil },
		"startOffset":         func(o *Options) { o.Reader.StartOffset = "middle" },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			o := readerOptions()
			change(o)
			_, err := NewReaderWithOptions(context.Background(), o)
			assert.Error(t, err)
		})
	}
}

func TestNewReaderWithOptions_AppliesPlugins(t *testing.T) {
	plugin := func(_ context.Context, c *kafka.ReaderConfig) error {
		c.HeartbeatInterval = 7 * time.Second
		return nil
	}
	r, err := NewReaderWithOptions(context.Background(), readerOptions(), plugin)
	require.NoError(t, err)
	t.Cleanup(func() { _ = r.Close() })
	assert.Equal(t, 7*time.Second, r.Config().HeartbeatInterval)

	failing := func(context.Context, *kafka.ReaderConfig) error { return errors.New("boom") }
	_, err = NewReaderWithOptions(context.Background(), readerOptions(), failing)
	assert.ErrorContains(t, err, "boom")
}

func TestNewConnWithOptions_Errors(t *testing.T) {
	o := &Options{Brokers: []string{"127.0.0.1:1"}, Conn: ConnOptions{Type: "leader", Network: "tcp", Timeout: time.Second}}

	_, err := NewConnWithOptions(context.Background(), o)
	assert.Error(t, err, "nothing listens on the broker address")

	o.Conn.Type = "follower"
	_, err = NewConnWithOptions(context.Background(), o)
	assert.ErrorContains(t, err, "follower")

	o.Conn.Type = "server"
	o.Brokers = nil
	_, err = NewConnWithOptions(context.Background(), o)
	assert.Error(t, err)
}

func TestNewConnWithOptions_PartitionDialsTheFirstBroker(t *testing.T) {
	o := &Options{Brokers: []string{"127.0.0.1:1"}, Conn: ConnOptions{Type: "partition", Network: "tcp", Timeout: time.Second}}
	_, err := NewConnWithOptions(context.Background(), o)
	assert.ErrorContains(t, err, "127.0.0.1:1")

	o.Brokers = []string{"no-port"}
	_, err = NewConnWithOptions(context.Background(), o)
	assert.ErrorContains(t, err, "no-port")
}
