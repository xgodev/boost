package main

import (
	"context"
	cfreecache "github.com/xgodev/boost/wrapper/cache/driver/contrib/coocood/freecache/v1"
	"github.com/xgodev/boost/wrapper/log/contrib/rs/zerolog/v1"
	"time"

	"github.com/coocood/freecache"
	"github.com/xgodev/boost/wrapper/cache"
	"github.com/xgodev/boost/wrapper/cache/codec/gob"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	ctx := context.Background()

	zerolog.NewLogger(zerolog.WithLevel("TRACE"))

	codec := gob.New[string]()

	fc1 := freecache.NewCache(1)
	fc2 := freecache.NewCache(1)

	drv1 := cfreecache.New(fc1, &cfreecache.Options{
		TTL: 10 * time.Minute,
	})

	drv2 := cfreecache.New(fc2, &cfreecache.Options{
		TTL: 10 * time.Minute,
	})

	v, _ := codec.Encode("value drv2")

	if err := drv2.Set(ctx, "key", v); err != nil {
		panic(err)
	}

	manager := cache.NewManager[string]("foo", codec, drv1, drv2)

	value, err := manager.GetOrSet(ctx, "key", func(ctx context.Context) (string, error) {
		return "value", nil
	})
	if err != nil {
		panic(err)
	}

	log.Infof("returned: %s", value)

	log.Infof("entries on drv1: %v", fc1.EntryCount())
	log.Infof("entries on drv2: %v", fc2.EntryCount())
}
