// Package receiver receives messages.
package receiver

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/consumer/internal/pkg/domain/model/entity"
)

// Service receives messages.
type Service interface {
	Receive(ctx context.Context, msg *entity.Message) error
}
