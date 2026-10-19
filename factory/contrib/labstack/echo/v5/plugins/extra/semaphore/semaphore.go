package semaphore

import (
	"context"
	e "github.com/labstack/echo/v5"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
	"golang.org/x/sync/semaphore"
)

// Register registers a new semaphore plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewSemaphore()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Semaphore represents semaphore plugin for echo server.
type Semaphore struct {
	options *Options
}

// NewSemaphoreWithOptions returns a new semaphore plugin with options.
func NewSemaphoreWithOptions(options *Options) *Semaphore {
	return &Semaphore{options: options}
}

// NewSemaphoreWithConfigPath returns a new semaphore plugin with options from config path.
func NewSemaphoreWithConfigPath(path string) (*Semaphore, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewSemaphoreWithOptions(o), nil
}

// NewSemaphore returns a new semaphore plugin with default options.
func NewSemaphore() (*Semaphore, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewSemaphoreWithOptions(o), nil
}

// Register registers this semaphore plugin for echo server.
func (i *Semaphore) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling semaphore middleware in echo")

	server.Use(semaphoreMiddleware(i.options.Limit))

	logger.Debug("semaphore middleware successfully enabled in echo")

	return nil
}

// semaphoreMiddleware answers 503 once limit requests are already in flight.
func semaphoreMiddleware(limit int64) e.MiddlewareFunc {
	sem := semaphore.NewWeighted(limit)

	return func(next e.HandlerFunc) e.HandlerFunc {
		return func(c *e.Context) error {
			if !sem.TryAcquire(1) {
				log.FromContext(c.Request().Context()).
					Errorf("the http server has reached the limit of %v simultaneous connections", limit)
				return e.ErrServiceUnavailable
			}
			defer sem.Release(1)

			return next(c)
		}
	}
}
