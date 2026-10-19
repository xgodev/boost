package kafka

import (
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Reader is the part of *kafka.Reader the subscriber uses.
type Reader interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
}

// Subscriber runs a handler on the messages of a reader.
type Subscriber[T any] struct {
	reader  Reader
	handler function.Handler[T]
	options *Options
}

// NewSubscriber returns a subscriber of reader.
func NewSubscriber[T any](reader Reader, handler function.Handler[T], options *Options) *Subscriber[T] {
	return &Subscriber[T]{reader: reader, handler: handler, options: options}
}

// Subscribe fetches messages until ctx ends or the reader is closed, then
// waits for the messages being handled. The messages of a partition are
// handled in order, one at a time; partitions run in parallel. A message is
// committed after its handler succeeds, or when it cannot be decoded or ran
// out of retries; one interrupted by the end of ctx stays uncommitted, to be
// delivered again. Any other fetch error is returned.
func (s *Subscriber[T]) Subscribe(ctx context.Context) error {

	logger := log.FromContext(ctx).WithTypeOf(*s)

	partitions := make(map[string]chan kafka.Message)
	var wg sync.WaitGroup
	defer func() {
		for _, messages := range partitions {
			close(messages)
		}
		wg.Wait()
	}()

	for {
		msg, err := s.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil || stderrors.Is(err, io.EOF) {
				logger.Infof("kafka subscriber stopped")
				return nil
			}
			return errors.Wrap(err, errors.Internalf("unable to fetch kafka message: %s", err))
		}

		key := fmt.Sprintf("%s-%d", msg.Topic, msg.Partition)
		messages, ok := partitions[key]
		if !ok {
			messages = make(chan kafka.Message, s.options.Buffer)
			partitions[key] = messages
			wg.Add(1)
			go func() {
				defer wg.Done()
				for m := range messages {
					s.process(ctx, m)
				}
			}()
		}

		select {
		case messages <- msg:
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *Subscriber[T]) process(ctx context.Context, msg kafka.Message) {

	logger := log.FromContext(ctx).WithTypeOf(*s)
	where := fmt.Sprintf("%s[%d]@%d", msg.Topic, msg.Partition, msg.Offset)

	in, err := decode(msg)
	if err != nil {
		logger.Errorf("skipping kafka message %s: %s", where, err)
		s.commit(ctx, msg, where)
		return
	}

	for retries := 0; ; retries++ {
		if ctx.Err() != nil {
			return
		}

		if _, err := s.handler(ctx, in); err == nil {
			break
		} else {
			logger.Errorf("handler failed on kafka message %s: %s", where, err)
		}

		if s.options.RetryLimit >= 0 && retries >= s.options.RetryLimit {
			logger.Errorf("giving up on kafka message %s after %d retries", where, retries)
			break
		}

		if !s.wait(ctx, retries+1) {
			return
		}
	}

	s.commit(ctx, msg, where)
}

func (s *Subscriber[T]) commit(ctx context.Context, msg kafka.Message, where string) {
	if err := s.reader.CommitMessages(ctx, msg); err != nil {
		log.FromContext(ctx).WithTypeOf(*s).Errorf("unable to commit kafka message %s: %s", where, err)
	}
}

// wait sleeps before retry n (1 for the first), doubling BackoffBase up to
// MaxBackoff. It returns false when ctx ends first.
func (s *Subscriber[T]) wait(ctx context.Context, n int) bool {
	if !s.options.Backoff {
		return ctx.Err() == nil
	}

	d := s.options.BackoffBase
	for i := 1; i < n && d < s.options.MaxBackoff; i++ {
		d *= 2
	}
	if d > s.options.MaxBackoff {
		d = s.options.MaxBackoff
	}

	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
	}
}
