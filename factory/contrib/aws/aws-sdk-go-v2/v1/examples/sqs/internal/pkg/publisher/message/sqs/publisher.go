// Package sqs publishes the messages on a queue, through the sqs client of
// boost's factory.
package sqs

import (
	"context"

	asqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sqs"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sqs/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sqs/internal/pkg/domain/publisher/message"
)

type publisher struct {
	client sqs.Client
	queue  string
}

// New loads the aws config and provides the publisher on the sqs client: the
// generator does not build aws clients, so the constructor calls the factory.
// @Provide(index=0)
func New(ctx context.Context, o *Options) (message.Publisher, error) {
	cfg, err := aws.NewConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &publisher{client: sqs.NewClient(asqs.NewFromConfig(cfg)), queue: o.Queue}, nil
}

func (p *publisher) Publish(ctx context.Context, msg *entity.Message) error {
	url, err := p.client.ResolveQueueUrl(ctx, p.queue)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, &asqs.SendMessageInput{QueueUrl: url, MessageBody: &msg.Text})
}
