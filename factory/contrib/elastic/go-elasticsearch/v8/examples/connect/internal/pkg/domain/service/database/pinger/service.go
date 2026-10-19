// Package pinger pings the database.
package pinger

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8/examples/connect/internal/pkg/domain/model/entity"
)

// Service pings the database.
type Service interface {
	Ping(ctx context.Context) (*entity.Database, error)
}
