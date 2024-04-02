package noop

import (
	"context"

	"github.com/xgodev/boost/wrapper/cache"
)

type noop struct {
}

func (c *noop) Set(ctx context.Context, key string, data []byte) error {
	return nil
}

func (c *noop) Del(ctx context.Context, key string) error {
	return nil
}

func (c *noop) Get(ctx context.Context, key string) (data []byte, err error) {
	return []byte{}, nil
}

func New() cache.Driver {
	return &noop{}
}
