package examplerepo

import (
	"github.com/xgodev/boost/wrapper/cache"
)

type plugin[R any] struct{}

func (l plugin[R]) Del(c *cache.Context[R], s string) error {
	return c.Del(s)
}

func (l plugin[R]) Get(c *cache.Context[R], s string) ([]byte, error) {
	return c.Get(s)
}

func (l plugin[R]) Set(c *cache.Context[R], s string, data []byte) error {
	return c.Set(s, data)
}

func New[R any]() cache.Plugin[R] {
	return &plugin[R]{}
}
