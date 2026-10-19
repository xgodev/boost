// Package sender sends messages.
package sender

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sqs/internal/pkg/domain/model/entity"
)

// Service sends a message.
type Service interface {
	Send(ctx context.Context, msg *entity.Message) error
}
