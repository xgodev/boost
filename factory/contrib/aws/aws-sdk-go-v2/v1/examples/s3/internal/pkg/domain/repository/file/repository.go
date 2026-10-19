// Package file reaches the files of the app.
package file

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/s3/internal/pkg/domain/model/entity"
)

// Repository reaches the files.
type Repository interface {
	// Find returns the file of the name, or a not found error.
	Find(ctx context.Context, name string) (*entity.File, error)
}
