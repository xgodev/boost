package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
	factory "github.com/xgodev/boost/factory/contrib/segmentio/kafka-go/v0"
	"github.com/xgodev/boost/model/errors"
)

// New returns the CmdFunc that runs the handler on the messages of the
// configured topics (the kafkago command), with the options read from config.
func New[T any]() function.CmdFunc[T] {
	return newCmd[T](NewOptions, factory.NewOptions)
}

// NewWithOptions returns the CmdFunc that runs the handler with the given
// adapter and reader options.
func NewWithOptions[T any](o *Options, fo *factory.Options) function.CmdFunc[T] {
	return newCmd[T](
		func() (*Options, error) { return o, nil },
		func() (*factory.Options, error) { return fo, nil },
	)
}

func newCmd[T any](options func() (*Options, error), factoryOptions func() (*factory.Options, error)) function.CmdFunc[T] {
	return func(fn function.Handler[T]) *co.Command {
		return &co.Command{
			Use:   "kafkago",
			Short: "kafkago",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				o, err := options()
				if err != nil {
					return err
				}
				fo, err := factoryOptions()
				if err != nil {
					return err
				}

				ctx := cmd.Context()
				r, err := newReader(ctx, fo, o.Topics)
				if err != nil {
					return err
				}
				// Close leaves the group, so its partitions move to the
				// other consumers right away.
				defer r.Close()

				return NewSubscriber[T](r, fn, o).Subscribe(ctx)
			},
		}
	}
}

// newReader returns the reader of the consumer group of fo on topics.
func newReader(ctx context.Context, fo *factory.Options, topics []string) (*kafka.Reader, error) {
	if fo.Reader.GroupID == "" {
		return nil, errors.NotValidf("the kafka-go adapter needs a consumer group: set boost.factory.kafkago.reader.groupId")
	}
	return factory.NewReaderWithOptions(ctx, fo, func(_ context.Context, c *kafka.ReaderConfig) error {
		c.Topic = ""
		c.GroupTopics = topics
		return nil
	})
}
