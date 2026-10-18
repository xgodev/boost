package function

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/extra/middleware"
	fxcontext "github.com/xgodev/boost/fx/modules/core/context"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

const (
	BSFunctionAdaptersGroupKey = "boostrap.function.adapters"
)

type params[T any] struct {
	fx.In
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Adapters   []function.CmdFunc[T] `group:"boostrap.function.adapters"`
}

var once sync.Once

// Module fx module for bootstrap function. The function runs from the app's
// OnStart and its context is cancelled on OnStop, so app.Run handles SIGTERM
// gracefully; when the function ends on its own the app shuts down, with
// exit code 1 if it failed.
func Module[T any](m []middleware.AnyErrorMiddleware[T]) fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			fxcontext.Module(),
			fx.Provide(
				func(ctx context.Context, p params[T]) (*function.Function[T], error) {
					return function.New[T](m...)
				}),
			fx.Invoke(register[T]),
		)
	})

	return options
}

func register[T any](ctx context.Context, p params[T], hdl function.Handler[T], fn *function.Function[T]) {
	runCtx, cancel := context.WithCancel(ctx)
	ran := make(chan struct{})
	var stopping atomic.Bool

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				defer close(ran)
				err := fn.Run(runCtx, hdl, p.Adapters...)
				stopped := stopping.Load()
				logger := log.FromContext(ctx)
				if err != nil && !(stopped && errors.Is(err, context.Canceled)) {
					logger.Errorf("function stopped with error: %v", err)
				}
				if stopped {
					return
				}
				code := 0
				if err != nil {
					code = 1
				}
				if err := p.Shutdowner.Shutdown(fx.ExitCode(code)); err != nil {
					logger.Error(err)
				}
			}()
			return nil
		},
		OnStop: func(stopCtx context.Context) error {
			stopping.Store(true)
			cancel()

			select {
			case <-ran:
				return nil
			case <-stopCtx.Done():
				return stopCtx.Err()
			}
		},
	})
}
