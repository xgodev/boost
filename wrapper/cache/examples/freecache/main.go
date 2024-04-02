package main

import (
	"context"
	"fmt"
	cfreecache "github.com/xgodev/boost/wrapper/cache/driver/contrib/coocood/freecache/v1"
	"time"

	"github.com/coocood/freecache"
	"github.com/xgodev/boost/wrapper/cache"
	"github.com/xgodev/boost/wrapper/cache/codec/gob"
)

func main() {

	fc := freecache.NewCache(1)

	drv := cfreecache.New(fc, &cfreecache.Options{
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
