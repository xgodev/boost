// Package local pings the database through the repository.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8/examples/connect/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8/examples/connect/internal/pkg/domain/repository/database"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8/examples/connect/internal/pkg/domain/service/database/pinger"
	"github.com/xgodev/boost/wrapper/log"
)

type service struct {
	repo database.Repository
}

// New provides the pinger on the database repository.
// @Provide(index=0)
func New(repo database.Repository) pinger.Service {
	return &service{repo: repo}
}

func (s *service) Ping(ctx context.Context) (*entity.Database, error) {
	db, err := s.repo.Ping(ctx)
	if err != nil {
		return nil, err
	}
	log.FromContext(ctx).Infof("%s status: %s", db.Name, db.Status)
	return db, nil
}
