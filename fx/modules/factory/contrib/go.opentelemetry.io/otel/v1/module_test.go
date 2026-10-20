package otel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1/oteltest"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var (
	// module holds Module() from its first call: the module is once-guarded
	// and yields nothing after that.
	module fx.Option
	// collector is where the providers the module starts export to.
	collector *oteltest.Collector
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	var err error
	collector, err = oteltest.NewCollector()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer collector.Close()
	if err := os.Setenv("BOOST_FACTORY_OTEL_ENDPOINT", collector.Endpoint()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	boost.Start()
	module = Module()
	return m.Run()
}

func TestModuleGraphIsComplete(t *testing.T) {
	require.NoError(t, fx.ValidateApp(module))
}

// Stopping the app must export what the providers still hold -- the metric
// export interval (60s) is far longer than the test, so only the stop can
// have sent the counter -- and shut the providers down.
func TestStopExportsWhatIsPendingAndShutsTheProvidersDown(t *testing.T) {
	ctx := context.Background()
	app := fxtest.New(t, module)
	app.RequireStart()

	tracer := otelboost.TracerProvider.Tracer("boost-test")
	counter, err := otelboost.MeterProvider.Meter("boost-test").Int64Counter("boost.test.calls")
	require.NoError(t, err)

	_, span := tracer.Start(ctx, "boost.test.before")
	require.True(t, span.IsRecording(), "the module starts the providers")
	span.End()
	counter.Add(ctx, 1)

	app.RequireStop()

	require.Contains(t, collector.Spans(), "boost.test.before")
	require.Contains(t, collector.Metrics(), "boost.test.calls")

	_, after := otelboost.TracerProvider.Tracer("boost-test").Start(ctx, "boost.test.after")
	require.False(t, after.IsRecording(), "the stop shut the providers down: a tracer taken after it is a noop")
}

// errorLogger keeps the messages logged at ERROR.
type errorLogger struct {
	log.Noop
	errors []string
}

func (l *errorLogger) FromContext(context.Context) log.Logger { return l }

func (l *errorLogger) Errorf(format string, args ...interface{}) {
	l.errors = append(l.errors, fmt.Sprintf(format, args...))
}

func TestShutdownLogsTheErrorInsteadOfFailingTheStop(t *testing.T) {
	previous := log.GetLogger()
	t.Cleanup(func() { log.Set(previous) })
	logger := &errorLogger{}
	log.Set(logger)

	shutdown(context.Background(), func(context.Context) error { return errors.New("collector unreachable") })
	shutdown(context.Background(), func(context.Context) error { return nil })

	require.Equal(t, []string{"opentelemetry: collector unreachable"}, logger.errors)
}
