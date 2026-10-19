// Package replier replies to the messages the function receives.
package replier

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/xgodev/boost/bootstrap/function/examples/cloudevents/http/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/bootstrap/function/examples/cloudevents/http/internal/pkg/domain/service/message/replier"
	"github.com/xgodev/boost/model/errors"
)

const (
	// RepliedType is the type of the event Handle replies with.
	RepliedType = "message.replied"
	// RepliesSubject is where the publisher middleware publishes the reply.
	RepliesSubject = "replies"
)

// Data is the data of the events in and out.
type Data struct {
	Message string `json:"message"`
}

// Handler replies to messages.
type Handler struct {
	svc replier.Service
}

// New provides the handler; the generated function door gets it from here.
// @Provide(index=0)
func New(svc replier.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle replies to the message of the event. The publisher middleware
// publishes the reply on subject replies; the logger middleware logs it.
// @Function(adapter=cloudevents, app=worker, middlewares=[publisher, logger])
func (h *Handler) Handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {
	var data Data
	if err := in.DataAs(&data); err != nil {
		return nil, errors.Wrap(err, errors.NotValidf("event %s has no message", in.ID()))
	}
	reply, err := h.svc.Reply(ctx, &entity.Message{Text: data.Message})
	if err != nil {
		return nil, err
	}
	out := cloudevents.NewEvent()
	out.SetID(uuid.NewString())
	out.SetType(RepliedType)
	out.SetSource("replier")
	out.SetSubject(RepliesSubject)
	if err := out.SetData(cloudevents.ApplicationJSON, Data{Message: reply.Text}); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error encoding the reply"))
	}
	return &out, nil
}
