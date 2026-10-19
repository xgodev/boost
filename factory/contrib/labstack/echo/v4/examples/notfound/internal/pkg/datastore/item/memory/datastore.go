// Package memory keeps the items in memory.
package memory

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/notfound/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/notfound/internal/pkg/domain/repository/item"
	"github.com/xgodev/boost/model/errors"
)

type datastore struct {
	items map[string]entity.Item
}

// New provides the items, with the item hello.
// @Provide(index=0)
func New() item.Repository {
	return &datastore{items: map[string]entity.Item{
		"hello": {ID: "hello", Name: "hello world"},
	}}
}

func (d *datastore) Get(_ context.Context, id string) (*entity.Item, error) {
	it, ok := d.items[id]
	if !ok {
		return nil, errors.NotFoundf("item %s", id)
	}
	return &it, nil
}
