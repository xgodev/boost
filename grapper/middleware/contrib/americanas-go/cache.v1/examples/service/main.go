package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/coocood/freecache"
	"github.com/xgodev/boost/cache"
	codec_cache_gob "github.com/xgodev/boost/cache/codec/gob"
	drv_cache_freecache "github.com/xgodev/boost/cache/driver/contrib/coocood/freecache.v1"
	mid_cache_log "github.com/xgodev/boost/cache/middleware/contrib/americanas-go/log.v1"
	"github.com/xgodev/boost/errors"
	"github.com/xgodev/boost/grapper"
	mid_grapper_cache "github.com/xgodev/boost/grapper/middleware/contrib/americanas-go/cache.v1"
	mid_grapper_fallback "github.com/xgodev/boost/grapper/middleware/native/fallback"
	"github.com/xgodev/boost/log"
	"github.com/xgodev/boost/log/contrib/rs/zerolog.v1"
)

type Result struct {
	Code string
}

type FooService struct {
	wrapper *grapper.AnyErrorWrapper[Result]
}

func NewFooService(wrapper *grapper.AnyErrorWrapper[Result]) *FooService {
	return &FooService{wrapper: wrapper}
}

func (s *FooService) FooMethod(ctx context.Context) (Result, error) {
	return s.wrapper.Exec(ctx, "1",
		func(ctx context.Context) (Result, error) {
			// business rule
			log.Infof("my business rule")
			rand.Seed(time.Now().UnixNano())
			if n := rand.Intn(100); n > 50 {
				return Result{Code: "SUCCESS"}, nil
			} else {
				return Result{}, errors.New("business error")
			}
		},
		func(ctx context.Context, r Result, err error) (Result, error) {
			// fallback rule
			if err != nil {
				log.Warnf("my fallback business rule")
				r.Code = "ERROR"
				return r, err
			}
			return r, err
		})
}

func main() {

	ctx := context.Background()

	zerolog.NewLogger(zerolog.WithLevel("TRACE"))

	var r Result
	var err error

	// CACHE
	fc := freecache.NewCache(1)
	drv := drv_cache_freecache.New(fc, &drv_cache_freecache.Options{
		TTL: 10 * time.Minute,
	})

	cachem := cache.NewManager[Result]("XPTO", codec_cache_gob.New[Result](), drv)
	cachem.Use(mid_cache_log.New[Result]())

	// GRAPPER
	middlewares := []grapper.AnyErrorMiddleware[Result]{
		mid_grapper_cache.NewAnyErrorMiddleware[Result](ctx, cachem, cache.SaveEmpty, cache.AsyncSave),
		mid_grapper_fallback.NewAnyErrorMiddleware[Result](),
	}

	wrapper := grapper.NewAnyErrorWrapper[Result](ctx, "XPTO", middlewares...)

	foo := NewFooService(wrapper)
	r, err = foo.FooMethod(ctx)
	if err != nil {
		log.Error(err)
	}

	log.Infof("CODE: %s", r.Code)
}
