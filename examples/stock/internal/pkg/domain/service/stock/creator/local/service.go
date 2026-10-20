// Package local creates stock in process, on the stock repository.
package local

import (
	"context"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/repository/stock"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/creator"
	"github.com/xgodev/boost/model/errors"
)

type service struct {
	repo stock.Repository
}

// New provides the stock creator, wrapped in the log and prometheus
// middlewares.
// @Provide(index=0)
// @Middleware(log, prometheus)
func New(repo stock.Repository) creator.Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, st *entity.Stock) error {
	if st.SKU == "" {
		return errors.NotValidf("sku is required")
	}
	if err := entity.ValidSKU(st.SKU); err != nil {
		return err
	}
	if st.Units < 0 {
		return errors.NotValidf("units of %s cannot be negative", st.SKU)
	}
	return s.repo.Create(ctx, st)
}
