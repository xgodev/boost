package sqs

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Client knows how to publish on sqs.
type Client interface {
	// Publish publishes message on sns.
	Publish(ctx context.Context, input *sqs.SendMessageInput) error

	// ResolveQueueUrl resolves the URL of the queue.
	ResolveQueueUrl(ctx context.Context, queueName string) (*string, error)
}

type sqsClient interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
}

// Client holds client and resource name.
type client struct {
	client sqsClient

	mu        sync.RWMutex // guards queueUrls: publishes resolve queues concurrently
	queueUrls map[string]*string
}

// NewClient returns a initialized client.
func NewClient(c *sqs.Client) Client {
	return &client{client: c, queueUrls: map[string]*string{}}
}

// Publish publishes input message to the configured sqs queue.
func (c *client) Publish(ctx context.Context, input *sqs.SendMessageInput) error {

	logger := log.FromContext(ctx).
		WithTypeOf(client{}). // not *c: that copies the lock.
		WithField("subject", input.QueueUrl)

	logger.Tracef("sending message to sqs")

	response, err := c.client.SendMessage(ctx, input)
	if err != nil {
		return errors.Wrap(err, errors.New("error sending message to sqs"))
	}

	logger.
		WithField("message_id", *response.MessageId).
		Debug("message sent to sqs")

	return nil
}

// ResolveQueueUrl resolves sqs queue url according to queueName.
func (c *client) ResolveQueueUrl(ctx context.Context, queueName string) (*string, error) {
	c.mu.RLock()
	queueUrl, ok := c.queueUrls[queueName]
	c.mu.RUnlock()
	if ok {
		return queueUrl, nil
	}

	result, err := c.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})

	if err != nil {
		return nil, errors.Wrap(err, errors.New("error resolving sqs queue url"))
	}

	if result == nil || result.QueueUrl == nil {
		return nil, errors.Errorf("sqs queue %s not found", queueName)
	}

	c.mu.Lock()
	c.queueUrls[queueName] = result.QueueUrl
	c.mu.Unlock()

	return result.QueueUrl, nil
}
