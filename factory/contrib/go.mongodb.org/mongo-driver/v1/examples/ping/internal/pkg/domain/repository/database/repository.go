// Package database reaches the database of the app.
package database

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1/examples/ping/internal/pkg/domain/model/entity"
)

// Repository reaches the database.
type Repository interface {
	// Ping asks the database whether it answers.
	Ping(ctx context.Context) (*entity.Database, error)
}
