package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coocood/freecache"
	"github.com/xgodev/boost/cache"
	"github.com/xgodev/boost/cache/codec/gob"
	driver "github.com/xgodev/boost/cache/driver/contrib/coocood/freecache.v1"
)

func main() {

	fc := freecache.NewCache(1)

	drv := driver.New(fc, &driver.Options{
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
