package log_test

import (
	"context"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus/v1"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus/v1/formatter/text"

	"github.com/xgodev/boost/wrapper/log"
)

func ExampleNewLogger() {
	log.SetGlobalLogger(logrus.NewLogger(logrus.WithFormatter(text.New(text.WithDisableTimestamp(true)))))
	log.WithField("main_field", "example")
	log.Info("main method.")
	// Output: level=info msg="main method."
}

func ExampleToContext() {
	bar := func(ctx context.Context) {
		logger := log.FromContext(ctx)

		logger = logger.WithField("bar_field", "example")
		logger.Infof("%s method.", "bar")
	}

	foo := func(ctx context.Context) {
		logger := log.FromContext(ctx)

		logger = logger.WithField("foo_field", "example")
		logger.Infof("%s method.", "foo")

		ctx = logger.ToContext(ctx)

		bar(ctx)
	}

	withoutContext := func() {
		log.Info("withoutContext method")
	}

	ctx := context.Background()
	log.SetGlobalLogger(logrus.NewLogger(
		logrus.WithFormatter(text.New(text.WithDisableTimestamp(true))),
	).WithField("main_field", "example"))

	ctx = log.ToContext(ctx)

	foo(ctx)

	withoutContext()

	// Output:
	// level=info msg="foo method." foo_field=example main_field=example
	// level=info msg="bar method." bar_field=example foo_field=example main_field=example
	// level=info msg="withoutContext method" main_field=example
}

func ExampleFromContext() {
	bar := func(ctx context.Context) {
		logger := log.FromContext(ctx)

		logger = logger.WithField("bar_field", "example")
		logger.Infof("%s method.", "bar")
	}

	foo := func(ctx context.Context) {
		logger := log.FromContext(ctx)

		logger = logger.WithField("foo_field", "example")
		logger.Infof("%s method.", "foo")

		ctx = logger.ToContext(ctx)

		bar(ctx)
	}

	withoutContext := func() {
		log.Info("withoutContext method")
	}

	ctx := context.Background()
	log.SetGlobalLogger(logrus.NewLogger(
		logrus.WithFormatter(text.New(text.WithDisableTimestamp(true))),
	).WithField("main_field", "example"))

	ctx = log.ToContext(ctx)

	foo(ctx)

	withoutContext()

	// Output:
	// level=info msg="foo method." foo_field=example main_field=example
	// level=info msg="bar method." bar_field=example foo_field=example main_field=example
	// level=info msg="withoutContext method" main_field=example
}
