package main

import (
	"context"
	mid_grapper_cache "github.com/xgodev/boost/middleware/plugins/local/wrapper/cache"
	cfreecache "github.com/xgodev/boost/wrapper/cache/driver/contrib/coocood/freecache/v1"
	mid_cache_log "github.com/xgodev/boost/wrapper/cache/plugins/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log/contrib/rs/zerolog/v1"
	"math/rand"
	"time"

	"github.com/coocood/freecache"
	"github.com/xgodev/boost/middleware"
	mid_grapper_fallback "github.com/xgodev/boost/middleware/plugins/native/fallback"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/cache"
	codec_cache_gob "github.com/xgodev/boost/wrapper/cache/codec/gob"
	"github.com/xgodev/boost/wrapper/log"
)

type Result struct {
	Code string
}

type FooService struct {
	wrapper *middleware.AnyErrorWrapper[Result]
}

func NewFooService(wrapper *middleware.AnyErrorWrapper[Result]) *FooService {
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
	drv := cfreecache.New(fc, &cfreecache.Options{
		TTL: 10 * time.Minute,
	})

	cachem := cache.NewManager[Result]("XPTO", codec_cache_gob.New[Result](), drv)
	cachem.Use(mid_cache_log.New[Result]())

	// GRAPPER
	middlewares := []middleware.AnyErrorMiddleware[Result]{
		mid_grapper_cache.NewAnyErrorMiddleware[Result](ctx, cachem, cache.SaveEmpty, cache.AsyncSave),
		mid_grapper_fallback.NewAnyErrorMiddleware[Result](),
	}

	wrapper := middleware.NewAnyErrorWrapper[Result](ctx, "XPTO", middlewares...)

	foo := NewFooService(wrapper)
	r, err = foo.FooMethod(ctx)
	if err != nil {
		log.Error(err)
	}

	log.Infof("CODE: %s", r.Code)
}
