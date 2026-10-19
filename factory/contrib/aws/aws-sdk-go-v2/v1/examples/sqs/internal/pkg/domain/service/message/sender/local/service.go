// Package local sends messages through the publisher.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sqs/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sqs/internal/pkg/domain/publisher/message"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sqs/internal/pkg/domain/service/message/sender"
	"github.com/xgodev/boost/model/errors"
)

type service struct {
	pub message.Publisher
}

// New provides the sender on the message publisher.
// @Provide(index=0)
func New(pub message.Publisher) sender.Service {
	return &service{pub: pub}
}

func (s *service) Send(ctx context.Context, msg *entity.Message) error {
	if msg.Text == "" {
		return errors.BadRequestf("the message has no text")
	}
	return s.pub.Publish(ctx, msg)
}
