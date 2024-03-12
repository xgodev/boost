package main

import (
	"context"
	"fmt"
	"github.com/xgodev/boost/cache/driver/contrib/coocood/freecache/v1"
	"time"

	"github.com/coocood/freecache"
	"github.com/xgodev/boost/cache"
	"github.com/xgodev/boost/cache/codec/gob"
)

func main() {

	fc := freecache.NewCache(1)

	drv := v1.New(fc, &v1.Options{
		TTL: 10 * time.Minute,
	})

	manager := cache.NewManager[string]("foo", gob.New[string](), drv)

	ctx := context.Background()

	if err := manager.Set(ctx, "key", "value"); err != nil {
		panic(err)
	}

	ok, value, err := manager.Get(ctx, "key")
	if err != nil {
		panic(err)
	}

	if !ok {
		fmt.Println("no key found")
	}

	fmt.Println(value)
}
