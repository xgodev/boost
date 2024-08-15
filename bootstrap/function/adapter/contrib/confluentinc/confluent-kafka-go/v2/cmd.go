package confluent

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
)

// New returns CmdFunc for cloudevents command.
func New[T any](consumer *kafka.Consumer) function.CmdFunc[T] {
	return func(fn function.Handler[T]) *co.Command {
		return &co.Command{
			Use:   "kafka_confluent",
			Short: "kafka_confluent",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				helper := NewHelper[T](consumer, fn)
				helper.Start()
				return nil
			},
		}
	}
}
