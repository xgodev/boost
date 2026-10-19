// Package echoer echoes the numbers the function receives.
package echoer

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/xgodev/boost/bootstrap/function/examples/middlewares/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/bootstrap/function/examples/middlewares/internal/pkg/domain/service/number/echoer"
	"github.com/xgodev/boost/model/errors"
)

// EchoedType is the type of the event Handle replies with.
const EchoedType = "number.echoed"

// Data is the data of the events in and out.
type Data struct {
	Number *int `json:"number"`
}

// Handler echoes numbers.
type Handler struct {
	svc echoer.Service
}

// New provides the handler; the generated function door gets it from here.
// @Provide(index=0)
func New(svc echoer.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle echoes the number of the event. An event without number panics, on
// purpose: the recovery middleware turns the panic into an error, the
// adapter answers it with 500 and the worker keeps running.
//
// prometheus runs first and recovery last, around the handler, so the logger
// and prometheus see the recovered panic as an error.
// @Function(adapter=cloudevents, app=worker, middlewares=[prometheus, logger, recovery])
func (h *Handler) Handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {
	var data Data
	if err := in.DataAs(&data); err != nil {
		return nil, errors.Wrap(err, errors.NotValidf("event %s has no number", in.ID()))
	}
	if data.Number == nil {
		panic("event " + in.ID() + " has no number")
	}
	number, err := h.svc.Echo(ctx, &entity.Number{Value: *data.Number})
	if err != nil {
		return nil, err
	}
	out := cloudevents.NewEvent()
	out.SetID(uuid.NewString())
	out.SetType(EchoedType)
	out.SetSource("echoer")
	if err := out.SetData(cloudevents.ApplicationJSON, Data{Number: &number.Value}); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error encoding the echo"))
	}
	return &out, nil
}
