// Package message publishes messages.
package message

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/sqs/internal/pkg/domain/model/entity"
)

// Publisher publishes a message.
type Publisher interface {
	Publish(ctx context.Context, msg *entity.Message) error
}
