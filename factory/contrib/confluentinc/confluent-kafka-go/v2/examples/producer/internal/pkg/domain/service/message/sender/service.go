// Package sender sends messages.
package sender

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/producer/internal/pkg/domain/model/entity"
)

// Service sends a message.
type Service interface {
	Send(ctx context.Context, msg *entity.Message) error
}
