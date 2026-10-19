// Package replier replies to messages.
package replier

import (
	"context"

	"github.com/xgodev/boost/bootstrap/function/examples/cloudevents/http/internal/pkg/domain/model/entity"
)

// Service replies to messages.
type Service interface {
	Reply(ctx context.Context, msg *entity.Message) (*entity.Message, error)
}
