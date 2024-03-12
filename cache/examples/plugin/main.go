package main

import (
	"context"
	cfreecache "github.com/xgodev/boost/cache/driver/contrib/coocood/freecache/v1"
	logger "github.com/xgodev/boost/cache/plugins/local/log"
	"github.com/xgodev/boost/log/contrib/rs/zerolog/v1"
	"time"

	"github.com/coocood/freecache"
	"github.com/xgodev/boost/cache"
	"github.com/xgodev/boost/cache/codec/gob"
	"github.com/xgodev/boost/log"
)

func main() {

	zerolog.NewLogger(zerolog.WithLevel("TRACE"))

	fc := freecache.NewCache(1)

	drv := cfreecache.New(fc, &cfreecache.Options{
		TTL: 10 * time.Minute,
	})

	manager := cache.NewManager[string]("foo", gob.New[string](), drv)
	manager.Use(logger.New[string]())

	ctx := context.Background()

	if err := manager.Set(ctx, "key", "value"); err != nil {
		panic(err)
	}

	ok, value, err := manager.Get(ctx, "key")
	if err != nil {
		panic(err)
	}

	if !ok {
		log.Infof("no key found")
	}

	value2, err := manager.GetOrSet(ctx, "key2", func(ctx context.Context) (string, error) {
		log.FromContext(ctx).Infof("executed get or set")
		return "get or set", nil
	}, cache.SaveEmpty)
	if err != nil {
		panic(err)
	}

	log.Infof(value)
	log.Infof(value2)
}
