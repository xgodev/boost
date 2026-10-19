// Package freecache keeps the readings in a cache manager over freecache,
// packed with msgpack.
package freecache

import (
	"context"
	"strconv"
	"time"

	ffreecache "github.com/xgodev/boost/factory/contrib/coocood/freecache/v1"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/cache"
	"github.com/xgodev/boost/wrapper/cache/codec/contrib/vmihailenco/msgpack/v5"
	dfreecache "github.com/xgodev/boost/wrapper/cache/driver/contrib/coocood/freecache/v1"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache-binary/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache-binary/internal/pkg/domain/repository/reading"
)

// packet is how the cache stores a reading: the sensor is the key, so the
// packet carries the rest.
type packet struct {
	Locid  uint16
	Tstamp uint32
	Temp   int16
}

type datastore struct {
	manager *cache.Manager[packet]
}

// New provides the readings on a freecache sized by
// boost.factory.freecache.cacheSize, where a reading lives 10 minutes. The
// manager encodes a packet with the msgpack codec.
// @Provide(index=0)
func New(ctx context.Context) (reading.Repository, error) {
	fc, err := ffreecache.NewCache(ctx)
	if err != nil {
		return nil, err
	}
	drv := dfreecache.New(fc, &dfreecache.Options{TTL: 10 * time.Minute})
	return &datastore{manager: cache.NewManager[packet]("readings", msgpack.New[packet](), drv)}, nil
}

func key(sensor uint32) string {
	return strconv.FormatUint(uint64(sensor), 10)
}

func (d *datastore) Save(ctx context.Context, r *entity.Reading) error {
	return d.manager.Set(ctx, key(r.Sensor), packet{Locid: r.Location, Tstamp: r.Timestamp, Temp: r.Temperature})
}

func (d *datastore) Get(ctx context.Context, sensor uint32) (*entity.Reading, error) {
	ok, p, err := d.manager.Get(ctx, key(sensor))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.NotFoundf("reading of sensor %d", sensor)
	}
	return &entity.Reading{Sensor: sensor, Location: p.Locid, Timestamp: p.Tstamp, Temperature: p.Temp}, nil
}
