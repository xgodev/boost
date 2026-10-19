// Package local greets the sender of the message.
package local

import (
	"context"

	"github.com/xgodev/boost/bootstrap/function/examples/nats/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/bootstrap/function/examples/nats/internal/pkg/domain/service/message/replier"
	"github.com/xgodev/boost/model/errors"
)

type service struct{}

// New provides the replier.
// @Provide(index=0)
func New() replier.Service {
	return &service{}
}

func (s *service) Reply(_ context.Context, msg *entity.Message) (*entity.Message, error) {
	if msg.Text == "" {
		return nil, errors.NotValidf("the message has no text")
	}
	return &entity.Message{Text: "hello, " + msg.Text}, nil
}
