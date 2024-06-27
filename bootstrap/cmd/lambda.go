package cmd

import (
	co "github.com/spf13/cobra"
	gsfx "github.com/xgodev/boost/bootstrap/fx"
	"github.com/xgodev/boost/bootstrap/lambda"
	"go.uber.org/fx"
)

// NewLambda returns CmdFunc for lambda command.
func NewLambda() CmdFunc {

	return func(options fx.Option) *co.Command {
		return &co.Command{
			Use:   "lambda",
			Short: "lambda",
			Long:  "",
			RunE: func(CmdFunc *co.Command, args []string) error {
				return gsfx.Run(lambda.HelperModule(options))
			},
		}
	}
}
