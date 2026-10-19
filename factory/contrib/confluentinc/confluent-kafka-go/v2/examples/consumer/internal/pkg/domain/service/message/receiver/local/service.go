// Package local receives the messages by logging them.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/consumer/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/consumer/internal/pkg/domain/service/message/receiver"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

type service struct{}

// New provides the receiver.
// @Provide(index=0)
func New() receiver.Service {
	return &service{}
}

func (s *service) Receive(ctx context.Context, msg *entity.Message) error {
	if msg.Text == "" {
		return errors.NotValidf("the message has no text")
	}
	log.FromContext(ctx).Infof("received %s", msg.Text)
	return nil
}
