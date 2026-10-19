package kafka

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
	factory "github.com/xgodev/boost/factory/contrib/segmentio/kafka-go/v0"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var (
	writerOnce sync.Once
	readerOnce sync.Once
)

// WriterModule provides the *kafka.Writer of boost.factory.kafkago. The
// writer is closed on app stop, which flushes pending messages.
//
// It returns an empty option after the first call.
func WriterModule() fx.Option {
	options := fx.Options()

	writerOnce.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle) (*kafka.Writer, error) {
					w, err := factory.NewWriter(ctx)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(w.Close))
					return w, nil
				},
			),
		)
	})

	return options
}

// ReaderModule provides the *kafka.Reader of boost.factory.kafkago. The
// reader is closed on app stop.
//
// It returns an empty option after the first call.
func ReaderModule() fx.Option {
	options := fx.Options()

	readerOnce.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle) (*kafka.Reader, error) {
					r, err := factory.NewReader(ctx)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(r.Close))
					return r, nil
				},
			),
		)
	})

	return options
}
