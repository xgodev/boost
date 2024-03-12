package grapper

import "context"

type AnyErrorContext[R any] struct {
	ctx   context.Context
	m     []AnyErrorMiddleware[R]
	name  string
	index int
	id    string
}

func (c *AnyErrorContext[R]) GetName() string {
	return c.name
}

func (c *AnyErrorContext[R]) GetContext() context.Context {
	return c.ctx
}

func (c *AnyErrorContext[R]) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *AnyErrorContext[R]) SetID(id string) {
	c.id = id
}

func (c *AnyErrorContext[R]) GetID() string {
	return c.id
}

func (c *AnyErrorContext[R]) Next(exec AnyErrorExecFunc[R], returnFunc AnyErrorReturnFunc[R]) (R, error) {
	if m := c.getNext(); m != nil {
		return m.Exec(c, exec, returnFunc)
	}
	return exec(c.GetContext())
}

func (c *AnyErrorContext[R]) hasNext() bool {
	if c.index < len(c.m) {
		return true
	}
	return false

}

func (c *AnyErrorContext[R]) getNext() AnyErrorMiddleware[R] {
	if c.hasNext() {
		m := c.m[c.index]
		c.index++
		return m
	}
	return nil
}

func NewAnyErrorContext[R any](name string, id string, m ...AnyErrorMiddleware[R]) *AnyErrorContext[R] {
	return &AnyErrorContext[R]{m: m, name: name, id: id}
}

type AnyContext[R any] struct {
	ctx   context.Context
	m     []AnyMiddleware[R]
	name  string
	index int
	id    string
}

func (c *AnyContext[R]) GetName() string {
	return c.name
}

func (c *AnyContext[R]) GetContext() context.Context {
	return c.ctx
}

func (c *AnyContext[R]) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *AnyContext[R]) SetID(id string) {
	c.id = id
}

func (c *AnyContext[R]) GetID() string {
	return c.id
}

func (c *AnyContext[R]) Next(exec AnyExecFunc[R], returnFunc AnyReturnFunc[R]) R {
	if m := c.getNext(); m != nil {
		return m.Exec(c, exec, returnFunc)
	}
	return exec(c.GetContext())
}

func (c *AnyContext[R]) hasNext() bool {
	if c.index < len(c.m) {
		return true
	}
	return false

}

func (c *AnyContext[R]) getNext() AnyMiddleware[R] {
	if c.hasNext() {
		m := c.m[c.index]
		c.index++
		return m
	}
	return nil
}

func NewAnyContext[R any](name string, id string, m ...AnyMiddleware[R]) *AnyContext[R] {
	return &AnyContext[R]{m: m, name: name, id: id}
}

type ErrorContext struct {
	ctx   context.Context
	m     []ErrorMiddleware
	name  string
	index int
	id    string
}

func (c *ErrorContext) GetName() string {
	return c.name
}

func (c *ErrorContext) GetContext() context.Context {
	return c.ctx
}

func (c *ErrorContext) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *ErrorContext) SetID(id string) {
	c.id = id
}

func (c *ErrorContext) GetID() string {
	return c.id
}

func (c *ErrorContext) Next(exec ErrorExecFunc, returnFunc ErrorReturnFunc) error {
	if m := c.getNext(); m != nil {
		return m.Exec(c, exec, returnFunc)
	}
	return exec(c.GetContext())
}

func (c *ErrorContext) hasNext() bool {
	if c.index < len(c.m) {
		return true
	}
	return false

}

func (c *ErrorContext) getNext() ErrorMiddleware {
	if c.hasNext() {
		m := c.m[c.index]
		c.index++
		return m
	}
	return nil
}

func NewErrorContext(name string, id string, m ...ErrorMiddleware) *ErrorContext {
	return &ErrorContext{m: m, name: name, id: id}
}
