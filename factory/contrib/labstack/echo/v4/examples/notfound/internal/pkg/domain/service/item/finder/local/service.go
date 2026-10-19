// Package local finds items in process, on the item repository.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/notfound/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/notfound/internal/pkg/domain/repository/item"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/notfound/internal/pkg/domain/service/item/finder"
)

type service struct {
	repo item.Repository
}

// New provides the finder.
// @Provide(index=0)
func New(repo item.Repository) finder.Service {
	return &service{repo: repo}
}

func (s *service) Find(ctx context.Context, id string) (*entity.Item, error) {
	return s.repo.Get(ctx, id)
}
