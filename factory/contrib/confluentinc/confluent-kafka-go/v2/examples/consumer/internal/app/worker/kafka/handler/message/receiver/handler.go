// Package receiver handles the records of the topic as messages.
package receiver

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/consumer/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/consumer/internal/pkg/domain/service/message/receiver"
)

// Handler receives the messages of the topic.
type Handler struct {
	svc receiver.Service
}

// New provides the handler; the generated function door gets it from here.
// @Provide(index=0)
func New(svc receiver.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle receives the value of the record as a message. An error makes the
// kafka adapter retry the record.
// @Function(adapter=kafka, app=worker)
func (h *Handler) Handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {
	return nil, h.svc.Receive(ctx, &entity.Message{Text: string(in.Data())})
}
