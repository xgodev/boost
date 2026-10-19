// Package database reaches the database of the app.
package database

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8/examples/connect/internal/pkg/domain/model/entity"
)

// Repository reaches the database.
type Repository interface {
	// Ping asks the database whether it answers.
	Ping(ctx context.Context) (*entity.Database, error)
}
