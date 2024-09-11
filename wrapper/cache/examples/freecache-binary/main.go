package main

import (
	"context"
	"fmt"
	"github.com/xgodev/boost/wrapper/cache/codec/contrib/vmihailenco/msgpack/v5"
	cfreecache "github.com/xgodev/boost/wrapper/cache/driver/contrib/coocood/freecache/v1"
	"time"

	"github.com/coocood/freecache"
	"github.com/xgodev/boost/wrapper/cache"
)

type packet struct {
	Sensid uint32
	Locid  uint16
	Tstamp uint32
	Temp   int16
}

func main() {

	fc := freecache.NewCache(1)

	drv := cfreecache.New(fc, &cfreecache.Options{
		TTL: 10 * time.Minute,
	})

	codec := msgpack.New[packet]()

	manager := cache.NewManager[packet]("foo", codec, drv)

	ctx := context.Background()

	data := packet{Sensid: 1, Locid: 1233, Tstamp: 123452123, Temp: 12}

	if err := manager.Set(ctx, "key", data); err != nil {
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
