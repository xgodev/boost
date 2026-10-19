// Package sender sends messages.
package sender

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0/examples/pubsub/internal/pkg/domain/model/entity"
)

// Service sends a message.
type Service interface {
	Send(ctx context.Context, msg *entity.Message) error
}
