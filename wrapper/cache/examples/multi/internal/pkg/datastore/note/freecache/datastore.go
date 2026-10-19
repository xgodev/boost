// Package freecache keeps the notes in a cache manager over two freecache
// layers.
package freecache

import (
	"context"
	"time"

	"github.com/coocood/freecache"
	ffreecache "github.com/xgodev/boost/factory/contrib/coocood/freecache/v1"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/cache"
	scodec "github.com/xgodev/boost/wrapper/cache/codec/string"
	dfreecache "github.com/xgodev/boost/wrapper/cache/driver/contrib/coocood/freecache/v1"
	"github.com/xgodev/boost/wrapper/cache/examples/multi/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/wrapper/cache/examples/multi/internal/pkg/domain/repository/note"
)

type datastore struct {
	near, far *freecache.Cache
	manager   *cache.Manager[string]
}

// New provides the notes on two freecaches sized by
// boost.factory.freecache.cacheSize: the near layer keeps a note 1 minute,
// the far one 10 minutes. The manager gets the drivers in the order it reads
// them, near first.
// @Provide(index=0)
func New(ctx context.Context) (note.Repository, error) {
	near, err := ffreecache.NewCache(ctx)
	if err != nil {
		return nil, err
	}
	far, err := ffreecache.NewCache(ctx)
	if err != nil {
		return nil, err
	}
	manager := cache.NewManager[string]("notes", scodec.New[string](),
		dfreecache.New(near, &dfreecache.Options{TTL: time.Minute}),
		dfreecache.New(far, &dfreecache.Options{TTL: 10 * time.Minute}),
	)
	return &datastore{near: near, far: far, manager: manager}, nil
}

// Save writes the note to every layer.
func (d *datastore) Save(ctx context.Context, n *entity.Note) error {
	return d.manager.Set(ctx, n.Key, n.Text)
}

// Get reads the note through GetOrSet, which copies a note found in the far
// layer to the near one. A key no layer has runs the func, whose NotFound
// GetOrSet returns without storing anything.
func (d *datastore) Get(ctx context.Context, key string) (*entity.Note, error) {
	text, err := d.manager.GetOrSet(ctx, key, func(context.Context) (string, error) {
		return "", errors.NotFoundf("note %s", key)
	})
	if err != nil {
		return nil, err
	}
	return &entity.Note{Key: key, Text: text}, nil
}
