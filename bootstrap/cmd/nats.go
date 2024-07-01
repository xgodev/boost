package cmd

import (
	co "github.com/spf13/cobra"
	gsfx "github.com/xgodev/boost/bootstrap/fx"
	"github.com/xgodev/boost/bootstrap/nats"
	"go.uber.org/fx"
)

// NewNats returns CmdFunc for nats command.
func NewNats() CmdFunc {
	return func(options fx.Option) *co.Command {
		return &co.Command{
			Use:   "nats",
			Short: "nats",
			Long:  "",
			RunE: func(CmdFunc *co.Command, args []string) error {
				return gsfx.Run(nats.HelperModule(options))
			},
		}
	}
}
