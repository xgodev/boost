// Package message publishes messages.
package message

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0/examples/pubsub/internal/pkg/domain/model/entity"
)

// Publisher publishes a message.
type Publisher interface {
	Publish(ctx context.Context, msg *entity.Message) error
}
