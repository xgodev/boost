// Package s3 keeps the files in a bucket, through the aws config of boost's
// factory.
package s3

import (
	"context"
	stderrors "errors"

	as3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/s3/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/s3/internal/pkg/domain/repository/file"
	"github.com/xgodev/boost/model/errors"
)

type repository struct {
	client *as3.Client
	bucket string
}

// New loads the aws config and provides the repository on an s3 client of it:
// the factory has no s3 client, so the constructor builds it from the config.
// @Provide(index=0)
func New(ctx context.Context, o *Options) (file.Repository, error) {
	cfg, err := aws.NewConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &repository{client: as3.NewFromConfig(cfg), bucket: o.Bucket}, nil
}

func (r *repository) Find(ctx context.Context, name string) (*entity.File, error) {
	head, err := r.client.HeadObject(ctx, &as3.HeadObjectInput{Bucket: &r.bucket, Key: &name})
	var notFound *types.NotFound
	if stderrors.As(err, &notFound) {
		return nil, errors.NotFoundf("file %s not found in bucket %s", name, r.bucket)
	}
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error reading file %s of bucket %s", name, r.bucket))
	}
	f := &entity.File{Name: name}
	if head.LastModified != nil {
		f.LastModified = *head.LastModified
	}
	return f, nil
}
