// Package freecache keeps the notes in a cache manager over freecache.
package freecache

import (
	"context"
	"time"

	ffreecache "github.com/xgodev/boost/factory/contrib/coocood/freecache/v1"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/cache"
	scodec "github.com/xgodev/boost/wrapper/cache/codec/string"
	dfreecache "github.com/xgodev/boost/wrapper/cache/driver/contrib/coocood/freecache/v1"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache/internal/pkg/domain/repository/note"
)

type datastore struct {
	manager *cache.Manager[string]
}

// New provides the notes on a freecache sized by
// boost.factory.freecache.cacheSize, where a note lives 10 minutes. The
// manager keeps the text of a note under its key, with the string codec.
// @Provide(index=0)
func New(ctx context.Context) (note.Repository, error) {
	fc, err := ffreecache.NewCache(ctx)
	if err != nil {
		return nil, err
	}
	drv := dfreecache.New(fc, &dfreecache.Options{TTL: 10 * time.Minute})
	return &datastore{manager: cache.NewManager[string]("notes", scodec.New[string](), drv)}, nil
}

func (d *datastore) Save(ctx context.Context, n *entity.Note) error {
	return d.manager.Set(ctx, n.Key, n.Text)
}

func (d *datastore) Get(ctx context.Context, key string) (*entity.Note, error) {
	ok, text, err := d.manager.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.NotFoundf("note %s", key)
	}
	return &entity.Note{Key: key, Text: text}, nil
}
