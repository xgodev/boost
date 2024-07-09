package nats

import (
	"github.com/nats-io/nats.go"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
)

// New returns CmdFunc for cloudevents command.
func New(conn *nats.Conn) function.CmdFunc {
	return func(fn function.Handler) *co.Command {
		return &co.Command{
			Use:   "nats",
			Short: "nats",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				helper := NewDefaultHelper(conn, fn)
				helper.Start()
				return nil
			},
		}
	}
}
