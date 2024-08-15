package confluent

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"
	"time"
)

// Helper assists in creating event handlers.
type Helper[T any] struct {
	handler  function.Handler[T]
	timeOut  time.Duration
	topics   []string
	consumer *kafka.Consumer
}

// NewHelperWithOptions returns a new Helper with options.
func NewHelperWithOptions[T any](consumer *kafka.Consumer, handler function.Handler[T], options *Options) *Helper[T] {

	return &Helper[T]{
		handler:  handler,
		timeOut:  options.TimeOut,
		topics:   options.Topics,
		consumer: consumer,
	}
}

// NewHelper returns a new Helper with default options.
func NewHelper[T any](consumer *kafka.Consumer, handler function.Handler[T]) *Helper[T] {

	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelperWithOptions(consumer, handler, opt)
}

func (h *Helper[T]) Start() {

	subscriber := NewSubscriber[T](h.consumer, h.handler, h.topics, h.timeOut)
	err := subscriber.Subscribe(context.Background())
	if err != nil {
		log.Error(err)
	}

}
