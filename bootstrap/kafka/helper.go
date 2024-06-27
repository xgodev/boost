package kafka

import (
	"context"
	"encoding/json"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/segmentio/kafka-go"
	"github.com/xgodev/boost/bootstrap/cloudevents"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"golang.org/x/sync/semaphore"
)

// Helper assists in creating event handlers.
type Helper struct {
	handler *cloudevents.HandlerWrapper
	options *Options
}

// NewHelper returns a new Helper with options.
func NewHelper(ctx context.Context, options *Options,
	handler *cloudevents.HandlerWrapper) *Helper {

	return &Helper{
		handler: handler,
		options: options,
	}
}

// NewDefaultHelper returns a new Helper with default options.
func NewDefaultHelper(ctx context.Context, handler *cloudevents.HandlerWrapper) *Helper {

	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelper(ctx, opt, handler)
}

func (h *Helper) Start() {

	for _, sub := range h.options.Subjects {
		sub := sub
		go h.subscribe(context.Background(), sub)
	}

	c := make(chan struct{})
	<-c
}

func (h *Helper) subscribe(ctx context.Context, topic string) {

	ctx = log.WithTypeOf(*h).
		WithField("topic", topic).
		WithField("groupId", h.options.GroupId).ToContext(ctx)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          h.options.Brokers,
		GroupID:          h.options.GroupId,
		Topic:            topic,
		Logger:           &Logger{},
		ErrorLogger:      &ErrorLogger{},
		QueueCapacity:    h.options.QueueCapacity,
		MinBytes:         h.options.MinBytes,
		MaxBytes:         h.options.MaxBytes,
		StartOffset:      h.options.StartOffset,
		ReadBatchTimeout: h.options.ReadBatchTimeout,
		MaxWait:          h.options.MaxWait,
		/*
			GroupTopics:            nil,
			Partition:              0,
			Dialer:                 nil,
			ReadLagInterval:        0,
			GroupBalancers:         nil,
			HeartbeatInterval:      0,
			CommitInterval:         0,
			PartitionWatchInterval: 0,
			WatchPartitionChanges:  false,
			SessionTimeout:         0,
			RebalanceTimeout:       0,
			JoinGroupBackoff:       0,
			RetentionTime:          0,
			ReadBackoffMin:         0,
			ReadBackoffMax:         0,
			IsolationLevel:         0,
			MaxAttempts:            0,
			OffsetOutOfRangeError:  false,
		*/
	})

	sem := semaphore.NewWeighted(int64(h.options.Concurrency))

	for {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Errorf(err.Error())
		}
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Errorf(err.Error())
			sem.Release(1)
			continue
		}
		go func(ctx context.Context, m kafka.Message) {
			defer sem.Release(1)
			ctx = log.FromContext(ctx).WithFields(
				map[string]interface{}{
					"kafka_partition": m.Partition,
					"kafka_topic":     m.Topic,
					"kafka_offset":    m.Offset,
				},
			).ToContext(ctx)
			h.handle(ctx, m)
		}(ctx, m)
	}

}

func (h *Helper) handle(ctx context.Context, msg kafka.Message) {

	logger := log.FromContext(ctx).WithTypeOf(*h)

	in := event.New()
	if err := json.Unmarshal(msg.Value, &in); err != nil {
		var data interface{}
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logger.Errorf("could not decode kafka record. %s", err.Error())
			return
		}

		err := in.SetData("", data)
		if err != nil {
			logger.Errorf("could set data from kafka record. %s", err.Error())
			return
		}

		// in.SetID(msg.)
	}

	var inouts []*cloudevents.InOut

	inouts = append(inouts, &cloudevents.InOut{In: &in})

	if err := h.handler.Process(ctx, inouts); err != nil {
		logger.Error(errors.ErrorStack(err))
	}

}

/*
type contextWithoutDeadline struct {
	ctx context.Context
}

func (*contextWithoutDeadline) Deadline() (time.Time, bool) { return time.Time{}, false }
func (*contextWithoutDeadline) Done() <-chan struct{}       { return nil }
func (*contextWithoutDeadline) Err() error                  { return nil }

func (l *contextWithoutDeadline) Value(key interface{}) interface{} {
	return l.ctx.Value(key)
}
*/
