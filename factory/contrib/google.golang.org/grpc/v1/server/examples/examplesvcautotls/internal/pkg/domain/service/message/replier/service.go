// Package replier replies to messages.
package replier

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/examples/examplesvcautotls/internal/pkg/domain/model/entity"
)

// Service replies to a message.
type Service interface {
	Reply(ctx context.Context, msg *entity.Message) (*entity.Message, error)
}
