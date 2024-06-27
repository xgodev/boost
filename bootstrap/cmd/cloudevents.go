package cmd

import (
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/cloudevents"
	gsfx "github.com/xgodev/boost/bootstrap/fx"
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
