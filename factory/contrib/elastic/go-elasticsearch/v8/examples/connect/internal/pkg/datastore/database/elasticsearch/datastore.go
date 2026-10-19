// Package elasticsearch reaches the cluster through the elasticsearch client
// of boost's factory (boost.factory.elasticsearch).
package elasticsearch

import (
	"context"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8"

	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8/examples/connect/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8/examples/connect/internal/pkg/domain/repository/database"
	"github.com/xgodev/boost/model/errors"
)

type repository struct {
	client *es.Client
}

// New connects the elasticsearch client and provides the repository on it:
// the generator does not build elasticsearch clients, so the constructor
// calls the factory.
// @Provide(index=0)
func New(ctx context.Context) (database.Repository, error) {
	client, err := elasticsearch.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &repository{client: client}, nil
}

func (r *repository) Ping(ctx context.Context) (*entity.Database, error) {
	res, err := r.client.Ping(r.client.Ping.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error pinging elasticsearch"))
	}
	defer res.Body.Close()
	return &entity.Database{Name: "elasticsearch", Status: res.Status()}, nil
}
