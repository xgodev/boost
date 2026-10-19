// Package ordercreated announces that an order was created.
package ordercreated

import (
	"context"

	"github.com/xgodev/boost/inject/examples/publisher/domain/model/entity"
)

// Publisher announces a created order.
type Publisher interface {
	// Publish announces order; a broker failure is returned as is.
	Publish(ctx context.Context, order entity.Order) error
}
