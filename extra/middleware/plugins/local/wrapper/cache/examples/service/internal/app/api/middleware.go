// Package api composes app api.
package api

import (
	"context"
	"time"

	"github.com/coocood/freecache"
	"github.com/xgodev/boost/extra/middleware"
	mcache "github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/cache"
	"github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/cache/examples/service/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/extra/middleware/plugins/native/fallback"
	ffreecache "github.com/xgodev/boost/factory/contrib/coocood/freecache/v1"
	"github.com/xgodev/boost/wrapper/cache"
	"github.com/xgodev/boost/wrapper/cache/codec/gob"
	dfreecache "github.com/xgodev/boost/wrapper/cache/driver/contrib/coocood/freecache/v1"
	clog "github.com/xgodev/boost/wrapper/cache/plugins/local/wrapper/log"
)

// Freecache is the in-memory cache the reports live in, sized by
// boost.factory.freecache.cacheSize.
// @Provide(index=0)
func Freecache(ctx context.Context) (*freecache.Cache, error) {
	return ffreecache.NewCache(ctx)
}

// Wrapper is what the finder calls run through: the cache middleware, which
// answers an id it has and stores what the rest of the chain returns
// without error, then the fallback middleware, which runs the fallback of the
// call. The cache keeps a report for 10 minutes and logs every access.
// @Provide(index=0)
func Wrapper(ctx context.Context, fc *freecache.Cache) *middleware.AnyErrorWrapper[*entity.Report] {
	manager := cache.NewManager[*entity.Report]("reports", gob.New[*entity.Report](),
		dfreecache.New(fc, &dfreecache.Options{TTL: 10 * time.Minute}))
	manager.Use(clog.New[*entity.Report]())
	return middleware.NewAnyErrorWrapper[*entity.Report](ctx, "reports",
		mcache.NewAnyErrorMiddleware[*entity.Report](ctx, manager, cache.SaveEmpty),
		fallback.NewAnyErrorMiddleware[*entity.Report](),
	)
}
