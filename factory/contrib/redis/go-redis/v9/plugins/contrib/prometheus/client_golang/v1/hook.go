package prometheus

import (
	"context"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/xgodev/boost/model/errors"
	"net"
	"time"
)

var (
	labelNames = []string{"command"}

	singleCommands = prom.NewHistogramVec(prom.HistogramOpts{
		Name:    "boost_factory_redis_single_commands",
		Help:    "Histogram of single Redis commands",
		Buckets: prom.DefBuckets,
	}, labelNames)

	pipelinedCommands = prom.NewCounterVec(prom.CounterOpts{
		Name: "boost_factory_redis_pipelined_commands",
		Help: "Number of pipelined Redis commands",
	}, labelNames)

	singleErrors = prom.NewCounterVec(prom.CounterOpts{
		Name: "boost_factory_redis_single_errors",
		Help: "Number of single Redis commands that have failed",
	}, labelNames)

	pipelinedErrors = prom.NewCounterVec(prom.CounterOpts{
		Name: "boost_factory_redis_pipelined_errors",
		Help: "Number of pipelined Redis commands that have failed",
	}, labelNames)
)

func init() {
	prom.MustRegister(singleCommands)
	prom.MustRegister(pipelinedCommands)
	prom.MustRegister(singleErrors)
	prom.MustRegister(pipelinedErrors)
}

type startKey struct{}

type Hook struct{}

func NewHook() redis.Hook {
	return &Hook{}
}

func (hook *Hook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (hook *Hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		duration := time.Since(start).Seconds()
		singleCommands.WithLabelValues(cmd.Name()).Observe(duration)

		if isActualErr(err) {
			singleErrors.WithLabelValues(cmd.Name()).Inc()
		}

		return err
	}
}

func (hook *Hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		err := next(ctx, cmds)

		for _, cmd := range cmds {
			pipelinedCommands.WithLabelValues(cmds[0].Name()).Inc()
			if isActualErr(cmd.Err()) {
				pipelinedErrors.WithLabelValues(cmd.Name()).Inc()
			}
		}

		return err
	}
}

func (hook *Hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey{}, time.Now()), nil
}

func (hook *Hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if start, ok := ctx.Value(startKey{}).(time.Time); ok {
		duration := time.Since(start).Seconds()
		singleCommands.WithLabelValues(cmd.Name()).Observe(duration)
	}

	if isActualErr(cmd.Err()) {
		singleErrors.WithLabelValues(cmd.Name()).Inc()
	}

	return nil
}

func (hook *Hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey{}, time.Now()), nil
}

func (hook *Hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if err := hook.AfterProcess(ctx, redis.NewCmd(ctx, "pipeline")); err != nil {
		return err
	}

	for _, cmd := range cmds {
		pipelinedCommands.WithLabelValues(cmd.Name()).Inc()

		if isActualErr(cmd.Err()) {
			pipelinedErrors.WithLabelValues(cmd.Name()).Inc()
		}
	}

	return nil
}

func isActualErr(err error) bool {
	return err != nil && !errors.Is(err, redis.Nil)
}
