package cmd

import (
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/faas/cloudevents"
	gsfx "github.com/xgodev/boost/faas/fx"
	"go.uber.org/fx"
)

// NewCloudEvents returns CmdFunc for cloudevents command.
func NewCloudEvents() CmdFunc {
	return func(options fx.Option) *co.Command {
		return &co.Command{
			Use:   "cloudevents",
			Short: "cloudevents",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				return gsfx.Run(cloudevents.HelperModule(options))
			},
		}
	}
}
