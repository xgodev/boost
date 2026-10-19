package lambda

import (
	"context"
	"sync"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"golang.org/x/sync/errgroup"
)

// Handler turns each record of a Lambda invocation into a CloudEvent and
// calls the function handler with it.
type Handler[T any] struct {
	fn      function.Handler[T]
	options *Options
}

// NewHandler returns a Handler that calls fn with the given options.
func NewHandler[T any](fn function.Handler[T], options *Options) *Handler[T] {
	return &Handler[T]{fn: fn, options: options}
}

// Handle is the Lambda entry point. Every record is handled, at most
// MaxWorkers at the same time; when any of them fails the invocation fails,
// so Lambda delivers the batch again. An error the errors registry ignores as
// success does not fail it.
func (h *Handler[T]) Handle(ctx context.Context, e Event) error {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return errors.Internalf("lambda context not found")
	}

	logger := log.FromContext(ctx).WithTypeOf(*h).WithField("awsRequestID", lc.AwsRequestID)
	ctx = logger.ToContext(ctx)

	if h.options.Skip {
		logger.Info("skipping event")
		return nil
	}

	ins, err := cloudEvents(e)
	if err != nil {
		return err
	}
	if len(ins) == 0 {
		logger.Warn("ignoring an invocation with no records")
		return nil
	}

	var mu sync.Mutex
	var failed []error

	var g errgroup.Group
	g.SetLimit(max(1, h.options.MaxWorkers))
	for _, in := range ins {
		g.Go(func() error {
			in.SetExtension("awsrequestid", lc.AwsRequestID)
			in.SetExtension("invokedfunctionarn", lc.InvokedFunctionArn)
			if _, err := h.fn(ctx, in); err != nil && !ignored(err) {
				logger.Errorf("handling event %s: %v", in.ID(), err)
				mu.Lock()
				failed = append(failed, err)
				mu.Unlock()
			}
			return nil
		})
	}
	_ = g.Wait() // intentional: the workers record their errors in failed and always return nil

	if len(failed) > 0 {
		return errors.Wrap(failed[0], errors.Internalf("%d of %d lambda records failed", len(failed), len(ins)))
	}
	return nil
}

// cloudEvents converts the invocation: one event per record, or the
// EventBridge event itself when it has no records.
func cloudEvents(e Event) ([]event.Event, error) {
	if len(e.Records) == 0 {
		if e.DetailType == "" {
			return nil, nil
		}
		in, err := fromEventBridge(e)
		if err != nil {
			return nil, err
		}
		return []event.Event{in}, nil
	}

	ins := make([]event.Event, 0, len(e.Records))
	for _, r := range e.Records {
		in, err := toCloudEvent(r)
		if err != nil {
			return nil, err
		}
		ins = append(ins, in)
	}
	return ins, nil
}

func ignored(err error) bool {
	pol, ok := errors.IgnoreOf(err)
	return ok && pol&errors.IgnoreAsSuccess != 0
}
