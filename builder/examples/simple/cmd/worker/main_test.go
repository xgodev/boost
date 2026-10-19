package main

import (
	"context"
	"os"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/bootstrap/function/middleware/logger"
	"github.com/xgodev/boost/bootstrap/function/middleware/recovery"
	"github.com/xgodev/boost/builder/examples/simple/gen/builder/worker"
	"github.com/xgodev/boost/builder/examples/simple/order"
	"github.com/xgodev/boost/builder/examples/simple/ordercreated"
	"github.com/xgodev/boost/extra/middleware"
	functionfx "github.com/xgodev/boost/fx/modules/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"go.uber.org/fx"
)

// TestMain boots boost: the function and its adapter read their config.
func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func event(t *testing.T, data any) cloudevents.Event {
	t.Helper()
	in := cloudevents.NewEvent()
	in.SetID("evt-1")
	in.SetType("order.created")
	in.SetSource("test")
	require.NoError(t, in.SetData(cloudevents.ApplicationJSON, data))
	return in
}

// TestWorker builds app worker from the generated Modules() and calls the
// function handler it provides, as the adapter does for each event.
func TestWorker(t *testing.T) {
	var (
		handle      function.Handler[*cloudevents.Event]
		svc         *order.Service
		middlewares []middleware.AnyErrorMiddleware[*cloudevents.Event]
	)
	app := fx.New(append(worker.Modules(), fx.NopLogger, fx.Populate(&handle, &svc),
		fx.Invoke(fx.Annotate(
			func(m []middleware.AnyErrorMiddleware[*cloudevents.Event]) { middlewares = m },
			fx.ParamTags(`name:"`+functionfx.MiddlewaresName+`"`),
		)),
	)...)
	require.NoError(t, app.Err())

	require.Len(t, middlewares, 2, "the function runs the middlewares of @Function")
	assert.IsType(t, &recovery.Recovery[*cloudevents.Event]{}, middlewares[0], "the outermost first")
	assert.IsType(t, &logger.Logger[*cloudevents.Event]{}, middlewares[1])

	out, err := handle(context.Background(), event(t, map[string]string{"item": "book"}))
	require.NoError(t, err)
	assert.Equal(t, ordercreated.RecordedType, out.Type())

	recorded, ok := svc.Get("1")
	require.True(t, ok, "the handler got the order service of the app")
	assert.Equal(t, order.Order{ID: "1", Item: "book"}, recorded)

	_, err = handle(context.Background(), event(t, map[string]string{}))
	assert.True(t, errors.IsNotValid(err), "got %v", err)
}
