package confluent

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"
)

// Subscriber represents a subscriber listener.
type Subscriber[T any] struct {
	consumer *kafka.Consumer
	handler  function.Handler[T]
	options  *Options
}

// NewSubscriber returns a subscriber listener.
func NewSubscriber[T any](consumer *kafka.Consumer, handler function.Handler[T], options *Options) *Subscriber[T] {
	return &Subscriber[T]{
		consumer: consumer,
		handler:  handler,
		options:  options,
	}
}

// Subscribe subscribes to specific topics and processes messages by partition.
func (l *Subscriber[T]) Subscribe(ctx context.Context) error {

	logger := log.FromContext(ctx)

	if err := l.consumer.SubscribeTopics(l.options.Topics, nil); err != nil {
		return err
	}

	// Update the map to use a string key composed of topic + partition
	partitions := make(map[string]chan *kafka.Message)

	for {
		msg, err := l.consumer.ReadMessage(l.options.TimeOut)
		if err != nil {
			if err.(kafka.Error).IsTimeout() {
				logger.Warnf("Consumer timeout: %v (%v)", err, msg)
				continue
			}
			logger.Errorf("Failed to read message: %v", err)
			continue
		}

		// Create a unique key using topic and partition
		topicPartitionKey := fmt.Sprintf("%s-%d", *msg.TopicPartition.Topic, msg.TopicPartition.Partition)

		// Check if a channel exists for the topic and partition
		if _, exists := partitions[topicPartitionKey]; !exists {
			partitions[topicPartitionKey] = make(chan *kafka.Message, l.options.MaxWorkers)

			// Process messages from each partition asynchronously if semaphore is used
			go l.processPartitionMessages(ctx, partitions[topicPartitionKey], topicPartitionKey)
		}

		partitions[topicPartitionKey] <- msg
	}
}

// processPartitionMessages processes messages for a specific topic and partition
func (l *Subscriber[T]) processPartitionMessages(ctx context.Context, messages chan *kafka.Message, topicPartitionKey string) {
	for msg := range messages {
		l.processMessage(ctx, msg, topicPartitionKey)
	}
}

// processMessage handles the actual message processing and retries
func (l *Subscriber[T]) processMessage(ctx context.Context, msg *kafka.Message, topicPartitionKey string) {
	logger := log.FromContext(ctx)
	retryCount := 0

	for {
		in := event.New()
		ce := false
		contentType := "application/json"

		if msg.Headers != nil {
			for _, h := range msg.Headers {
				switch h.Key {
				case "content-type":
					in.SetDataContentType(string(h.Value))
					contentType = string(h.Value)
				case "ce_specversion":
					in.SetSpecVersion(string(h.Value))
					ce = true
				case "ce_id":
					in.SetID(string(h.Value))
					ce = true
				case "ce_source":
					in.SetSource(string(h.Value))
					ce = true
				case "ce_type":
					in.SetType(string(h.Value))
					ce = true
				case "ce_time":
					if t, err := time.Parse(time.RFC3339, string(h.Value)); err == nil {
						in.SetTime(t)
					}
					ce = true
				case "ce_subject":
					in.SetSubject(string(h.Value))
					ce = true
				default:
					in.SetExtension(h.Key, string(h.Value))
				}
			}
		}

		if !ce {
			in.SetID(uuid.NewString())
			in.SetSource(fmt.Sprintf("kafka://%s/%v", *msg.TopicPartition.Topic, msg.TopicPartition.Partition))
			in.SetType("kafka.message")
			in.SetTime(time.Now())
		}

		if err := in.SetData(contentType, msg.Value); err != nil {
			logger.Warnf("could not set data from kafka record. %s", err.Error())
			continue
		}

		_, err := l.handler(ctx, in)
		if err != nil {
			logger.Errorf("Handler error in topic-partition %s: %v", topicPartitionKey, err)
			retryCount++

			// Check if retry limit is reached
			if l.options.RetryLimit != -1 && retryCount >= l.options.RetryLimit {
				logger.Errorf("Max retry limit reached for topic-partition %s. Message will not be retried.", topicPartitionKey)
				break
			}

			// Apply backoff if enabled
			if l.options.Backoff {
				l.applyBackoff(retryCount)
			}

			continue // Retry message processing
		}

		// Commit the message if manual commit is enabled
		if l.options.ManualCommit {
			if _, err := l.consumer.CommitMessage(msg); err != nil {
				logger.Errorf("Failed to commit message from topic-partition %s: %v", topicPartitionKey, err)
				continue // Retry on commit failure
			}
		}

		logger.Infof("Message from topic-partition %s successfully processed and committed: %s", topicPartitionKey, string(msg.Value))
		break // Exit loop on success
	}
}

// applyBackoff applies an exponential backoff strategy with a configurable base and max
func (l *Subscriber[T]) applyBackoff(retryCount int) {
	// Use exponential backoff based on the retry count
	backoffTime := time.Duration(math.Pow(2, float64(retryCount))) * l.options.BackoffBase

	// Cap the backoff time to the configured max
	if backoffTime > l.options.MaxBackoff {
		backoffTime = l.options.MaxBackoff
	}
	time.Sleep(backoffTime)
}
