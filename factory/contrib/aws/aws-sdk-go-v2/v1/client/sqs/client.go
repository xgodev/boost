package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
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
	client    sqsClient
	queueUrls map[string]*string
}

// NewClient returns a initialized client.
func NewClient(c *sqs.Client) Client {
	return &client{c, map[string]*string{}}
}

// Publish publishes input message to the configured sqs queue.
func (c *client) Publish(ctx context.Context, input *sqs.SendMessageInput) error {

	logger := log.FromContext(ctx).
		WithTypeOf(*c).
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
	if queueUrl, ok := c.queueUrls[queueName]; ok {
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

	c.queueUrls[queueName] = result.QueueUrl

	return result.QueueUrl, nil
}
