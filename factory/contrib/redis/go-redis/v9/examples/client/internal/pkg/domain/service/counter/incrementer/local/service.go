// Package local increments the counters kept by the repository.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9/examples/client/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9/examples/client/internal/pkg/domain/repository/counter"
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9/examples/client/internal/pkg/domain/service/counter/incrementer"
)

type service struct {
	repo counter.Repository
}

// New provides the incrementer on the counter repository.
// @Provide(index=0)
func New(repo counter.Repository) incrementer.Service {
	return &service{repo: repo}
}

func (s *service) Increment(ctx context.Context, name string) (*entity.Counter, error) {
	count, err := s.repo.Increment(ctx, name)
	if err != nil {
		return nil, err
	}
	return &entity.Counter{Name: name, Count: count}, nil
}
