// Package local replies hello world to every message.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvcautotls/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvcautotls/internal/pkg/domain/service/message/replier"
	"github.com/xgodev/boost/wrapper/log"
)

type service struct{}

// New provides the replier.
// @Provide(index=0)
func New() replier.Service {
	return &service{}
}

func (s *service) Reply(ctx context.Context, msg *entity.Message) (*entity.Message, error) {
	log.FromContext(ctx).Infof("message received: %s", msg.Text)
	return &entity.Message{Text: "hello world"}, nil
}
