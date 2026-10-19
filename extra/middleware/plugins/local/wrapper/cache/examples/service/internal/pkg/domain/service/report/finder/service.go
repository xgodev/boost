// Package finder finds reports.
package finder

import (
	"context"

	"github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/cache/examples/service/internal/pkg/domain/model/entity"
)

// Service finds reports.
type Service interface {
	Find(ctx context.Context, id string) (*entity.Report, error)
}
