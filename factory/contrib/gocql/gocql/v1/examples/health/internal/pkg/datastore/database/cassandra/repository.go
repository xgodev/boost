// Package cassandra reaches cassandra through a session of boost's factory
// (boost.factory.gocql).
package cassandra

import (
	"context"

	cql "github.com/gocql/gocql"
	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v1"
	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v1/examples/health/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v1/examples/health/internal/pkg/domain/repository/database"
	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v1/plugins/local/extra/health"
	"github.com/xgodev/boost/model/errors"
)

type repository struct {
	session *cql.Session
}

// New opens the session with the health plugin and provides the repository on it: the
// generator does not build cassandra sessions, so the constructor calls the
// factory.
// @Provide(index=0)
func New(ctx context.Context) (database.Repository, error) {
	session, err := gocql.NewSession(ctx, health.Register)
	if err != nil {
		return nil, err
	}
	return &repository{session: session}, nil
}

func (r *repository) Ping(ctx context.Context) (*entity.Database, error) {
	if err := r.session.Query("SELECT now() FROM system.local").WithContext(ctx).Exec(); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error querying cassandra"))
	}
	return &entity.Database{Name: "cassandra", Status: "ok"}, nil
}
