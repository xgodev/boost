package cmd

import (
	co "github.com/spf13/cobra"
	gsfx "github.com/xgodev/boost/bootstrap/fx"
	"github.com/xgodev/boost/bootstrap/pubsub"
	"go.uber.org/fx"
)

// NewPubSub returns CmdFunc for pubsub command.
func NewPubSub() CmdFunc {
	return func(options fx.Option) *co.Command {
		return &co.Command{
			Use:   "pubsub",
			Short: "pubsub",
			Long:  "",
			RunE: func(CmdFunc *co.Command, args []string) error {
				return gsfx.Run(pubsub.HelperModule(options))
			},
		}
	}
}
