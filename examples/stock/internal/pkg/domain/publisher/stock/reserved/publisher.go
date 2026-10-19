// Package reserved announces that a unit of an item was reserved.
package reserved

import (
	"context"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
)

// Publisher announces a reservation.
type Publisher interface {
	// Publish announces the stock left after a reservation.
	Publish(ctx context.Context, s *entity.Stock) error
}
