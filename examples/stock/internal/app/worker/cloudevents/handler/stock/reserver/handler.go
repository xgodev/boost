// Package reserver reserves the item of each order.created event in app
// worker.
package reserver

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/reserver"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Handler reserves the stock of orders.
type Handler struct {
	svc reserver.Service
}

// New provides the handler; the generated function door gets it from here.
// @Provide(index=0)
func New(svc reserver.Service) *Handler {
	return &Handler{svc: svc}
}

// order is the data of an order.created event.
type order struct {
	SKU string `json:"sku"`
}

// Handle reserves one unit of the sku of the order. The reservation is
// announced by the reserver, so the function replies with no event. The
// adapter is worker.function.adapter: cloudevents (HTTP, the default) or
// nats.
// @Function(adapter=[cloudevents, nats], app=worker, middlewares=[recovery, logger])
func (h *Handler) Handle(ctx context.Context, in cloudevents.Event) ([]*cloudevents.Event, error) {
	var o order
	if err := in.DataAs(&o); err != nil || o.SKU == "" {
		return nil, errors.NotValidf("event %s has no sku", in.ID())
	}
	left, err := h.svc.Reserve(ctx, o.SKU)
	if err != nil {
		return nil, err
	}
	log.FromContext(ctx).Infof("reserved one %s for event %s, %d left", left.SKU, in.ID(), left.Units)
	return nil, nil
}
