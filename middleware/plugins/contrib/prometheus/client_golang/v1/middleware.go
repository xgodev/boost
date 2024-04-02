package prometheus

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/xgodev/boost/middleware"
	"github.com/xgodev/boost/model/errors"
)

const (
	opsTotalName        = "ops_total"
	opsTotalErroredName = "ops_total_errored"
	opsDurationName     = "ops_duration_seconds"
)

var (
	buckets = []float64{
		0.0005,
		0.001, // 1ms
		0.002,
		0.005,
		0.01, // 10ms
		0.02,
		0.05,
		0.1, // 100 ms
		0.2,
		0.5,
		1.0, // 1s
		2.0,
		5.0,
		10.0, // 10s
		15.0,
		20.0,
		30.0,
	}

	opsTotalProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: opsTotalName,
		Help: "The total number of processed operations",
	}, []string{"name"})

	opsTotalErrored = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: opsTotalErroredName,
		Help: "The total number of processed operations with error",
	}, []string{"name", "error"})

	opsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    opsDurationName,
		Help:    "Spend time by processing a operation",
		Buckets: buckets,
	}, []string{"name"})
)

type anyErrorMiddleware[R any] struct {
}

func (c *anyErrorMiddleware[R]) Exec(ctx *middleware.AnyErrorContext[R], exec middleware.AnyErrorExecFunc[R], returnFunc middleware.AnyErrorReturnFunc[R]) (r R, err error) {
	opsTotalProcessed.WithLabelValues(ctx.GetName()).Inc()
	timer := prometheus.NewTimer(opsDuration.WithLabelValues(ctx.GetName()))
	defer timer.ObserveDuration()
	r, err = ctx.Next(exec, returnFunc)
	if err != nil {
		opsTotalErrored.WithLabelValues(ctx.GetName(), errString(err)).Inc()
	}
	return r, err
}

func NewAnyErrorMiddleware[R any](ctx context.Context) middleware.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{}
}

type anyMiddleware[R any] struct {
}

func (c *anyMiddleware[R]) Exec(ctx *middleware.AnyContext[R], exec middleware.AnyExecFunc[R], returnFunc middleware.AnyReturnFunc[R]) (r R) {
	opsTotalProcessed.WithLabelValues(ctx.GetName()).Inc()
	timer := prometheus.NewTimer(opsDuration.WithLabelValues(ctx.GetName()))
	defer timer.ObserveDuration()
	return ctx.Next(exec, returnFunc)
}

func NewAnyMiddleware[R any](ctx context.Context) middleware.AnyMiddleware[R] {
	return &anyMiddleware[R]{}
}

type errorMiddleware struct {
}

func (c *errorMiddleware) Exec(ctx *middleware.ErrorContext, exec middleware.ErrorExecFunc, returnFunc middleware.ErrorReturnFunc) (err error) {
	opsTotalProcessed.WithLabelValues(ctx.GetName()).Inc()
	timer := prometheus.NewTimer(opsDuration.WithLabelValues(ctx.GetName()))
	defer timer.ObserveDuration()
	if err = ctx.Next(exec, returnFunc); err != nil {
		opsTotalErrored.WithLabelValues(ctx.GetName(), errString(err)).Inc()
	}
	return err

}

func NewErrorMiddleware(ctx context.Context) middleware.ErrorMiddleware {
	return &errorMiddleware{}
}

func errString(err error) string {
	switch {
	case errors.IsNotFound(err):
		return "NOT_FOUND"
	case errors.IsMethodNotAllowed(err):
		return "NOT_ALLOWED"
	case errors.IsNotValid(err) || errors.IsBadRequest(err):
		return "NOT_VALID"
	case errors.IsServiceUnavailable(err):
		return "SERVICE_UNAVAILABLE"
	case errors.IsConflict(err) || errors.IsAlreadyExists(err):
		return "CONFLICT"
	case errors.IsNotImplemented(err) || errors.IsNotProvisioned(err):
		return "NOT_IMPLEMENTED"
	case errors.IsUnauthorized(err):
		return "UNAUTHORIZED"
	case errors.IsForbidden(err):
		return "FORBIDDEN"
	case errors.IsNotSupported(err) || errors.IsNotAssigned(err):
		return "NOT_SUPPORTED"
	case errors.IsInternal(err):
		return "INTERNAL_ERROR"
	default:
		if _, ok := err.(validator.ValidationErrors); ok {
			return "UNPROCESSABLE_ENTITY"
		}
		return "GENERIC_ERROR"
	}
}
