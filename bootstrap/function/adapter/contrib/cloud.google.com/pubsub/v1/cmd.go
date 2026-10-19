package pubsub

import (
	"cloud.google.com/go/pubsub"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
)

// New returns the CmdFunc of the gcp_pubsub command: it handles the messages
// of the configured subscriptions until the run context is done.
func New[T any](client *pubsub.Client) function.CmdFunc[T] {
	return func(fn function.Handler[T]) *co.Command {
		return &co.Command{
			Use:   "gcp_pubsub",
			Short: "gcp_pubsub",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				return NewHelper[T](client, fn).Run(cmd.Context())
			},
		}
	}
}
