package nats

import (
	"github.com/nats-io/nats.go"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
)

// New returns CmdFunc for cloudevents command.
func New[T any](conn *nats.Conn) function.CmdFunc[T] {
	return func(fn function.Handler[T]) *co.Command {
		return &co.Command{
			Use:   "cloudevents_nats",
			Short: "cloudevents_nats",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				helper := NewHelper[T](conn, fn)
				helper.Start()
				return nil
			},
		}
	}
}
