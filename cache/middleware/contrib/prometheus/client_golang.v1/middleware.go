package prometheus

import (
	"github.com/xgodev/boost/cache"
)

type middleware[R any] struct{}

func (l middleware[R]) Del(c *cache.Context[R], s string) error {
	return c.Del(s)
}

func (l middleware[R]) Get(c *cache.Context[R], s string) ([]byte, error) {
	return c.Get(s)
}

func (l middleware[R]) Set(c *cache.Context[R], s string, data []byte) error {
	return c.Set(s, data)
}

func New[R any]() cache.Middleware[R] {
	return &middleware[R]{}
}
