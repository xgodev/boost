// Package boost announces reservations on boost's wrapper/publisher as
// CloudEvents. The driver is boost.wrapper.publisher.driver.name: noop (the
// default) or nats, the only broker driver generated for this app.
package boost

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/publisher/stock/reserved"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/publisher"
)

// Type is the CloudEvents type of the reservation event.
const Type = "stock.reserved"

type boostPublisher struct {
	pub     *publisher.Publisher
	subject string
}

// New provides the reservation publisher on the publisher builder
// generates.
// @Provide(index=0)
// @Inject(index=0, drivers=[nats])
func New(pub *publisher.Publisher, o *Options) reserved.Publisher {
	return &boostPublisher{pub: pub, subject: o.Subject}
}

// event is the data of the reservation event.
type event struct {
	SKU   string `json:"sku"`
	Units int64  `json:"units"`
}

func (p *boostPublisher) Publish(ctx context.Context, s *entity.Stock) error {
	out := cloudevents.NewEvent()
	out.SetID(uuid.NewString())
	out.SetType(Type)
	out.SetSource("stock")
	out.SetSubject(p.subject)
	if err := out.SetData(cloudevents.ApplicationJSON, event{SKU: s.SKU, Units: s.Units}); err != nil {
		return errors.Wrap(err, errors.Internalf("error encoding the reservation of %s", s.SKU))
	}
	return p.pub.Publish(ctx, []*cloudevents.Event{&out})
}
