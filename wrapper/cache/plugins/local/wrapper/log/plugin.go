package log

import (
	"github.com/xgodev/boost/wrapper/cache"
	"github.com/xgodev/boost/wrapper/log"
)

type plugin[R any] struct{}

func (l plugin[R]) Del(c *cache.Context[R], s string) error {
	logger := l.logger(c, s)
	logger.Tracef("executing Del method")
	defer logger.Tracef("executed Del method")
	return c.Del(s)
}

func (l plugin[R]) Get(c *cache.Context[R], s string) ([]byte, error) {
	logger := l.logger(c, s)
	logger.Tracef("executing Get method")
	defer logger.Tracef("executed Get method")
	return c.Get(s)
}

func (l plugin[R]) Set(c *cache.Context[R], s string, data []byte) error {
	logger := l.logger(c, s)
	logger.Tracef("executing Set method")
	defer logger.Tracef("executed Set method")
	return c.Set(s, data)
}

func (l plugin[R]) logger(c *cache.Context[R], k string) log.Logger {
	return log.FromContext(c.GetContext()).
		WithField("cache_key", k).
		WithField("cache_name", c.GetName())
}

func New[R any]() cache.Plugin[R] {
	return &plugin[R]{}
}
