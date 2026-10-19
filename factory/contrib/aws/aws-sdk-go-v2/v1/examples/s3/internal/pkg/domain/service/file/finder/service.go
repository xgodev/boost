// Package finder finds files.
package finder

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/s3/internal/pkg/domain/model/entity"
)

// Service finds files.
type Service interface {
	Find(ctx context.Context, name string) (*entity.File, error)
}
