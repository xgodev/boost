package newrelic

import (
	"context"

	v2 "github.com/cloudevents/sdk-go/v2"
	nr "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/xgodev/boost/bootstrap/cloudevents"
	newrelic "github.com/xgodev/boost/factory/contrib/newrelic/go-agent/v3"
)

// NewRelic represents a newrelic agent middleware for events.
type NewRelic struct {
	cloudevents.UnimplementedMiddleware
}

// NewNewRelic creates a newrelic agent middleware.
func NewNewRelic() cloudevents.Middleware {
	if !IsEnabled() {
		return nil
	}
	return &NewRelic{}
}

// BeforeAll starts a newrelic transaction before processing all input event handlers.
// The transaction started is passed via context.
func (m *NewRelic) BeforeAll(ctx context.Context, inout []*cloudevents.InOut) (context.Context, error) {

	txn := newrelic.Application().StartTransaction(TxName())

	c := nr.NewContext(ctx, txn)

	return c, nil
}

// Before enables the newrelic transacation for use in multiple goroutines to be used by the handler.
func (m *NewRelic) Before(parentCtx context.Context, in *v2.Event) (context.Context, error) {

	txn := nr.FromContext(parentCtx).NewGoroutine()

	ctx := nr.NewContext(parentCtx, txn)

	return ctx, nil
}

// After checks if the handler has returned any error and notifies via newrelic agent.
func (m *NewRelic) After(parentCtx context.Context, in v2.Event, out *v2.Event, err error) (context.Context, error) {

	txn := nr.FromContext(parentCtx)

	if err != nil {
		if txn != nil {
			txn.NoticeError(err)
		}
	}

	return parentCtx, nil
}

// Close finishes the newrelic transaction.
func (m *NewRelic) Close(ctx context.Context) error {
	txn := nr.FromContext(ctx)
	defer txn.End()

	return nil
}
