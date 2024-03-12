package cache

import "context"

type Context[R any] struct {
	ctx    context.Context
	m      []Middleware[R]
	driver Driver
	name   string
	index  int
}

func (c *Context[R]) GetDriver() Driver {
	return c.driver
}

func (c *Context[R]) GetName() string {
	return c.name
}

func (c *Context[R]) GetContext() context.Context {
	return c.ctx
}

func (c *Context[R]) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *Context[R]) Del(key string) error {
	if m := c.getNext(); m != nil {
		return m.Del(c, key)
	}
	return c.driver.Del(c.GetContext(), key)
}

func (c *Context[R]) Get(key string) ([]byte, error) {
	if m := c.getNext(); m != nil {
		return m.Get(c, key)
	}
	return c.driver.Get(c.GetContext(), key)
}

func (c *Context[R]) Set(key string, data []byte) error {
	if m := c.getNext(); m != nil {
		return m.Set(c, key, data)
	}
	return c.driver.Set(c.GetContext(), key, data)
}

func (c *Context[R]) hasNext() bool {
	if c.index < len(c.m) {
		return true
	}
	return false

}
func (c *Context[R]) getNext() Middleware[R] {
	if c.hasNext() {
		m := c.m[c.index]
		c.index++
		return m
	}
	return nil
}

func NewContext[R any](name string, d Driver, m ...Middleware[R]) *Context[R] {
	return &Context[R]{m: m, name: name, driver: d}
}
