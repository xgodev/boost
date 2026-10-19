// Package local keeps notes in process, on the note repository.
package local

import (
	"context"

	"github.com/xgodev/boost/wrapper/cache/examples/freecache/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache/internal/pkg/domain/repository/note"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache/internal/pkg/domain/service/note/keeper"
)

type service struct {
	repo note.Repository
}

// New provides the keeper.
// @Provide(index=0)
func New(repo note.Repository) keeper.Service {
	return &service{repo: repo}
}

func (s *service) Save(ctx context.Context, n *entity.Note) error {
	return s.repo.Save(ctx, n)
}

func (s *service) Find(ctx context.Context, key string) (*entity.Note, error) {
	return s.repo.Get(ctx, key)
}
