// Package mongodb reaches mongo through the connection the generator builds
// on boost's factory (boost.factory.mongo).
package mongodb

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2/examples/ping/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2/examples/ping/internal/pkg/domain/repository/database"
	"github.com/xgodev/boost/model/errors"
)

type repository struct {
	conn *mongo.Conn
}

// New provides the repository on the mongo connection of the generator.
// @Provide(index=0)
// @Inject(index=0)
func New(conn *mongo.Conn) database.Repository {
	return &repository{conn: conn}
}

func (r *repository) Ping(ctx context.Context) (*entity.Database, error) {
	if err := r.conn.Client.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error pinging mongo"))
	}
	return &entity.Database{Name: "mongo", Status: "ok"}, nil
}
