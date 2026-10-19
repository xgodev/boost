// Package local keeps readings in process, on the reading repository.
package local

import (
	"context"

	"github.com/xgodev/boost/wrapper/cache/examples/freecache-binary/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache-binary/internal/pkg/domain/repository/reading"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache-binary/internal/pkg/domain/service/reading/keeper"
)

type service struct {
	repo reading.Repository
}

// New provides the keeper.
// @Provide(index=0)
func New(repo reading.Repository) keeper.Service {
	return &service{repo: repo}
}

func (s *service) Save(ctx context.Context, r *entity.Reading) error {
	return s.repo.Save(ctx, r)
}

func (s *service) Find(ctx context.Context, sensor uint32) (*entity.Reading, error) {
	return s.repo.Get(ctx, sensor)
}
