package job

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/xgodev/boost/bootstrap/job"
	fxcontext "github.com/xgodev/boost/fx/modules/core/context"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

type params struct {
	fx.In
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Handler    job.Handler
}

var once sync.Once

// Module fx module that runs the job.Handler the app provides once. It runs
// from the app's OnStart and shuts the app down when it returns: exit code 0
// on nil, 1 on an error, logged at ERROR. OnStop (SIGTERM) cancels its
// context and waits for it until the stop timeout; a handler that had not
// finished fails the stop, so the app exits 1.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			fxcontext.Module(),
			fx.Invoke(register),
		)
	})

	return options
}

func register(ctx context.Context, p params) {
	runCtx, cancel := context.WithCancel(ctx)
	ran := make(chan struct{})
	var err error
	// reported is set by whoever reports how the job ended: the handler
	// returning on its own, or OnStop
	var reported atomic.Bool

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				defer close(ran)
				err = p.Handler(runCtx)
				cancel()
				if !reported.CompareAndSwap(false, true) {
					return // the app is stopping: OnStop reports it
				}
				logger := log.FromContext(ctx)
				code := 0
				if err != nil {
					logger.Errorf("job failed: %v", err)
					code = 1
				}
				if err := p.Shutdowner.Shutdown(fx.ExitCode(code)); err != nil {
					logger.Error(err)
				}
			}()
			return nil
		},
		OnStop: func(stopCtx context.Context) error {
			if !reported.CompareAndSwap(false, true) {
				return nil // the handler returned and shut the app down
			}
			cancel()

			select {
			case <-ran:
				if err != nil {
					return fmt.Errorf("job stopped before finishing: %w", err)
				}
				return nil
			case <-stopCtx.Done():
				return fmt.Errorf("job did not return before the stop timeout: %w", stopCtx.Err())
			}
		},
	})
}
