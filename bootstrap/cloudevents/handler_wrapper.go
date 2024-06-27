package cloudevents

import (
	"context"
	"reflect"

	v2 "github.com/cloudevents/sdk-go/v2"
	cloudevents "github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/utils/strings"
	"github.com/xgodev/boost/wrapper/log"
)

// HandlerWrapper can be used to process events wrapped in middleware.
type HandlerWrapper struct {
	handler     cloudevents.Handler
	middlewares []Middleware
	options     *HandlerWrapperOptions
}

// NewHandlerWrapper creates a new handler with options and wrapped in middleware.
func NewHandlerWrapper(handler cloudevents.Handler, options *HandlerWrapperOptions, middlewares ...Middleware) *HandlerWrapper {

	if middlewares == nil {
		middlewares = []Middleware{}
	}

	return &HandlerWrapper{handler: handler, middlewares: middlewares, options: options}
}

// NewDefaultHandlerWrapper creates a new handler wrapped in middleware.
func NewDefaultHandlerWrapper(handler cloudevents.Handler, middlewares ...Middleware) *HandlerWrapper {

	opt, err := DefaultHandlerWrapperOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHandlerWrapper(handler, opt, middlewares...)
}

func (h *HandlerWrapper) closeAll(parentCtx context.Context) error {

	logger := log.FromContext(parentCtx).WithTypeOf(*h)

	for _, middleware := range h.middlewares {

		logger.Tracef("executing event middleware %s.Close()", reflect.TypeOf(middleware).String())

		var err error
		err = middleware.Close(parentCtx)
		if err != nil {
			err = errors.Wrap(err, errors.Internalf("an error happened when calling Close() method in %s middleware",
				reflect.TypeOf(middleware).String()))
			return err
		}
	}

	return nil
}

func (h *HandlerWrapper) afterAll(parentCtx context.Context, inouts []*InOut) error {

	logger := log.FromContext(parentCtx).WithTypeOf(*h)

	for _, middleware := range h.middlewares {

		logger.Tracef("executing event middleware %s.AfterAll()", reflect.TypeOf(middleware).String())

		var err error
		parentCtx, err = middleware.AfterAll(parentCtx, inouts)
		if err != nil {
			err = errors.Wrap(err,
				errors.Internalf("an error happened when calling AfterAll() method in %s middleware",
					reflect.TypeOf(middleware).String()))
			return err
		}
	}

	return nil
}

func (h *HandlerWrapper) beforeAll(parentCtx context.Context, inouts []*InOut) (context.Context, error) {

	logger := log.FromContext(parentCtx).WithTypeOf(*h)

	for _, middleware := range h.middlewares {

		logger.Tracef("executing event middleware %s.BeforeAll()", reflect.TypeOf(middleware).String())

		var err error
		parentCtx, err = middleware.BeforeAll(parentCtx, inouts)
		if err != nil {
			err = errors.Wrap(err,
				errors.Internalf("an error happened when calling BeforeAll() method in %s middleware",
					reflect.TypeOf(middleware).String()))
			return parentCtx, err
		}
	}

	return parentCtx, nil
}

// Process processes events.
func (h *HandlerWrapper) Process(parentCtx context.Context, inouts []*InOut) (err error) {

	logger := log.FromContext(parentCtx).WithTypeOf(*h)

	parentCtx, err = h.beforeAll(parentCtx, inouts)
	if err != nil {
		logger.Error(errors.ErrorStack(err))
		return err
	}

	h.handleAll(parentCtx, inouts)

	err = h.afterAll(parentCtx, inouts)
	if err != nil {
		logger.Error(errors.ErrorStack(err))
		return err
	}

	err = h.closeAll(parentCtx)
	if err != nil {
		logger.Error(errors.ErrorStack(err))
		return err
	}

	for _, inout := range inouts {
		if inout.Err != nil {
			return err
		}
	}

	return nil
}

func (h *HandlerWrapper) handleAll(parentCtx context.Context, inouts []*InOut) {

	for _, inout := range inouts {
		logger := log.FromContext(parentCtx)

		in := inout.In
		if in == nil {
			logger.Warn("discarding inout.In == nil")
			continue
		}

		logger = logger.
			WithField("event.id", in.ID()).
			WithField("event.parentId", in.Extensions()["parentId"]).
			WithField("event.source", in.Source()).
			WithField("event.type", in.Type())

		if inout.Err != nil {
			logger.WithField("cause", inout.Err.Error()).Warn("discarding message due to error")
			inout.Err = nil
			continue
		}

		hasToDiscardEvent := strings.SliceContains(h.options.IDsToDiscard, in.ID())
		if hasToDiscardEvent {
			logger.Warn("discarding event due to feature flag")
			continue
		}

		var err error
		ctx := logger.ToContext(parentCtx)

		ctx, err = h.before(ctx, h.middlewares, in)
		if err != nil {
			logger.WithField("cause", err.Error()).Warn("could not execute h.before")
		}

		var out *v2.Event
		out, err = h.handler(ctx, *in)
		if err != nil {
			inout.Err = errors.Wrap(err, errors.Internalf("unable process event"))
		}
		inout.Out = out
		inout.Context = ctx

		_, err = h.after(ctx, h.middlewares, *in, out, inout.Err)
		if err != nil {
			logger.WithField("cause", err.Error()).Warn("could not execute h.after")
		}
	}
}

func (h *HandlerWrapper) before(ctx context.Context, middlewares []Middleware, in *v2.Event) (context.Context, error) {

	logger := log.FromContext(ctx).WithTypeOf(*h)

	var err error

	for _, middleware := range middlewares {

		logger.Tracef("executing event middleware %s.Before()", reflect.TypeOf(middleware).String())

		ctx, err = middleware.Before(ctx, in)
		if err != nil {
			return ctx, errors.Wrap(err,
				errors.Internalf("an error happened when calling Before() method in %s middleware",
					reflect.TypeOf(middleware).String()))
		}
	}

	return ctx, nil
}

func (h *HandlerWrapper) after(ctx context.Context, middlewares []Middleware, in v2.Event, out *v2.Event,
	err error) (context.Context, error) {

	logger := log.FromContext(ctx).WithTypeOf(*h)

	var er error

	for _, middleware := range middlewares {

		logger.Tracef("executing event middleware %s.After()", reflect.TypeOf(middleware).String())

		ctx, er = middleware.After(ctx, in, out, err)
		if er != nil {
			return ctx, errors.Wrap(err,
				errors.Internalf("an error happened when calling After() method in %s middleware",
					reflect.TypeOf(middleware).String()))
		}
	}

	return ctx, nil
}
