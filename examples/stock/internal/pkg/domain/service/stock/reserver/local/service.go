// Package local reserves stock in process: it takes the unit from the
// repository and announces the reservation.
package local

import (
	"context"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/publisher/stock/reserved"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/repository/stock"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/reserver"
	"github.com/xgodev/boost/model/errors"
)

type service struct {
	repo stock.Repository
	pub  reserved.Publisher
}

// New provides the default stock reserver, wrapped in the log and
// prometheus middlewares.
// @Provide(index=0)
// @Default
// @Middleware(log, prometheus)
func New(repo stock.Repository, pub reserved.Publisher) reserver.Service {
	return &service{repo: repo, pub: pub}
}

func (s *service) Reserve(ctx context.Context, sku string) (*entity.Stock, error) {
	if err := entity.ValidSKU(sku); err != nil {
		return nil, err
	}
	left, err := s.repo.Take(ctx, sku)
	if err != nil {
		return nil, err
	}
	if err := s.pub.Publish(ctx, left); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error announcing the reservation of %s", sku))
	}
	return left, nil
}
