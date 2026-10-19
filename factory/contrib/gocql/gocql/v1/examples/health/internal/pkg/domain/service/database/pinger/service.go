// Package pinger pings the database.
package pinger

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v1/examples/health/internal/pkg/domain/model/entity"
)

// Service pings the database.
type Service interface {
	Ping(ctx context.Context) (*entity.Database, error)
}
