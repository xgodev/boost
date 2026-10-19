// Package gocloud sends the messages to a topic, through boost's gocloud
// factory.
package gocloud

import (
	"context"

	bpubsub "github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0"
	"github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0/examples/pubsub/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0/examples/pubsub/internal/pkg/domain/publisher/message"
	"github.com/xgodev/boost/model/errors"
	"gocloud.dev/pubsub"
)

type publisher struct {
	topic *pubsub.Topic
}

// New opens the topic of the factory and provides the publisher on it: the
// generator does not build gocloud topics, so the constructor calls the
// factory.
// @Provide(index=0)
func New(ctx context.Context) (message.Publisher, error) {
	topic, err := bpubsub.NewTopic(ctx)
	if err != nil {
		return nil, err
	}
	return &publisher{topic: topic}, nil
}

func (p *publisher) Publish(ctx context.Context, msg *entity.Message) error {
	if err := p.topic.Send(ctx, &pubsub.Message{Body: []byte(msg.Text)}); err != nil {
		return errors.Wrap(err, errors.Internalf("error sending the message"))
	}
	return nil
}
