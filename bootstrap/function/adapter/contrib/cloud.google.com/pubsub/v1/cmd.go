package pubsub

import (
	"cloud.google.com/go/pubsub"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
)

// New returns CmdFunc for cloudevents command.
func New[T any](client *pubsub.Client) function.CmdFunc[T] {
	return func(fn function.Handler[T]) *co.Command {
		return &co.Command{
			Use:   "gcp_pubsub",
			Short: "gcp_pubsub",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				helper := NewHelper[T](client, fn)
				helper.Start()
				return nil
			},
		}
	}
}
