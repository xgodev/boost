// Package counter is a custom middleware: it counts the calls that run
// through it.
package counter

import (
	"sync/atomic"

	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/wrapper/log"
)

// Counter counts the calls that run through it.
type Counter[R any] struct {
	calls atomic.Int64
}

// New returns a counter at zero.
func New[R any]() *Counter[R] {
	return &Counter[R]{}
}

// Exec counts the call and runs the next middleware, or the call itself.
func (c *Counter[R]) Exec(ctx *middleware.AnyErrorContext[R], exec middleware.AnyErrorExecFunc[R], fallback middleware.AnyErrorReturnFunc[R]) (R, error) {
	n := c.calls.Add(1)
	log.FromContext(ctx.GetContext()).Debugf("call %d of %s", n, ctx.GetName())
	return ctx.Next(exec, fallback)
}

// Calls is how many calls ran through the counter.
func (c *Counter[R]) Calls() int64 {
	return c.calls.Load()
}
