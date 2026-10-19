// Package local builds the reports it finds, through the middleware wrapper
// the app provides.
package local

import (
	"context"
	"strconv"
	"sync/atomic"

	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/cache/examples/service/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/cache/examples/service/internal/pkg/domain/service/report/finder"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

type service struct {
	wrapper *middleware.AnyErrorWrapper[*entity.Report]
	builds  atomic.Int64
}

// New provides the finder on the wrapper its calls run through.
// @Provide(index=0)
func New(wrapper *middleware.AnyErrorWrapper[*entity.Report]) finder.Service {
	return &service{wrapper: wrapper}
}

// Find runs the build of the report through the wrapper, with id as the
// cache key: the build only runs when the cache has no report for id.
func (s *service) Find(ctx context.Context, id string) (*entity.Report, error) {
	return s.wrapper.Exec(ctx, id,
		func(ctx context.Context) (*entity.Report, error) {
			if _, err := strconv.Atoi(id); err != nil {
				return nil, errors.NotValidf("report id %q is not a number", id)
			}
			return &entity.Report{ID: id, Build: s.builds.Add(1)}, nil
		},
		// the fallback, for a failed build: it logs and keeps the error, so
		// the cache middleware does not store the call
		func(ctx context.Context, r *entity.Report, err error) (*entity.Report, error) {
			if err != nil {
				log.FromContext(ctx).Warnf("report %s not built: %v", id, err)
			}
			return r, err
		})
}
