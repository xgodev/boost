// Package local finds stock in process, on the stock repository.
package local

import (
	"context"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/repository/stock"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/finder"
)

type service struct {
	repo stock.Repository
}

// New provides the stock finder, wrapped in the log and prometheus
// middlewares.
// @Provide(index=0)
// @Middleware(log, prometheus)
func New(repo stock.Repository) finder.Service {
	return &service{repo: repo}
}

func (s *service) Find(ctx context.Context, sku string) (*entity.Stock, error) {
	return s.repo.Get(ctx, sku)
}
