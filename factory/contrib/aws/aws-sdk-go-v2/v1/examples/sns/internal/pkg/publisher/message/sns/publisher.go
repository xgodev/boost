// Package sns publishes the messages on a topic, through the sns client of
// boost's factory.
package sns

import (
	"context"

	asns "github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sns"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sns/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sns/internal/pkg/domain/publisher/message"
)

type publisher struct {
	client sns.Client
	topic  string
}

// New loads the aws config and provides the publisher on the sns client: the
// generator does not build aws clients, so the constructor calls the factory.
// @Provide(index=0)
func New(ctx context.Context, o *Options) (message.Publisher, error) {
	cfg, err := aws.NewConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &publisher{client: sns.NewClient(asns.NewFromConfig(cfg)), topic: o.Topic}, nil
}

func (p *publisher) Publish(ctx context.Context, msg *entity.Message) error {
	return p.client.Publish(ctx, &asns.PublishInput{TopicArn: &p.topic, Message: &msg.Text})
}
