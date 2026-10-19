// Package local finds the files through the repository.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/s3/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/s3/internal/pkg/domain/repository/file"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/s3/internal/pkg/domain/service/file/finder"
	"github.com/xgodev/boost/wrapper/log"
)

type service struct {
	repo file.Repository
}

// New provides the finder on the file repository.
// @Provide(index=0)
func New(repo file.Repository) finder.Service {
	return &service{repo: repo}
}

func (s *service) Find(ctx context.Context, name string) (*entity.File, error) {
	f, err := s.repo.Find(ctx, name)
	if err != nil {
		return nil, err
	}
	log.FromContext(ctx).Infof("file %s last modified at %s", f.Name, f.LastModified)
	return f, nil
}
