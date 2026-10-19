package nats

import (
	"github.com/nats-io/nats.go"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
)

// New returns the CmdFunc of the nats command: it handles the messages of
// the configured subjects until the run context is done.
func New[T any](conn *nats.Conn) function.CmdFunc[T] {
	return func(fn function.Handler[T]) *co.Command {
		return &co.Command{
			Use:   "nats",
			Short: "nats",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				return NewHelper[T](conn, fn).Run(cmd.Context())
			},
		}
	}
}
