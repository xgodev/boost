package cmd

import (
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/azure"
	gsfx "github.com/xgodev/boost/bootstrap/fx"
	"go.uber.org/fx"
)

// NewAzure returns CmdFunc for azure functions command.
func NewAzure() CmdFunc {
	return func(options fx.Option) *co.Command {
		return &co.Command{
			Use:   "azure",
			Short: "azure",
			Long:  "",
			RunE: func(CmdFunc *co.Command, args []string) error {
				return gsfx.Run(azure.HelperModule(options))
			},
		}
	}
}
