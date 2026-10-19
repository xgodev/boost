// Package pinger pings the database.
package pinger

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1/examples/ping/internal/pkg/domain/model/entity"
)

// Service pings the database.
type Service interface {
	Ping(ctx context.Context) (*entity.Database, error)
}
