// Package local reserves stock in process, on the stock repository.
package local

import (
	"context"

	"github.com/acme/shop/internal/pkg/domain/repository/stock"
	"github.com/acme/shop/internal/pkg/domain/service/reservestock"
)

type service struct {
	repo stock.Repository
}

// New provides the service, wrapped in the log and prometheus middlewares.
// @Provide(index=0)
// @Middleware(log, prometheus)
func New(repo stock.Repository) reservestock.Service {
	return &service{repo: repo}
}

func (s *service) Reserve(ctx context.Context, sku string) error {
	return s.repo.Take(ctx, sku)
}
