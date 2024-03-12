package cmd

import (
	co "github.com/spf13/cobra"
	gsfx "github.com/xgodev/boost/faas/fx"
	"github.com/xgodev/boost/faas/kafka"
	"go.uber.org/fx"
)

// NewKafka returns CmdFunc for kafka command.
func NewKafka() CmdFunc {
	return func(options fx.Option) *co.Command {
		return &co.Command{
			Use:   "kafka",
			Short: "kafka",
			Long:  "",
			RunE: func(CmdFunc *co.Command, args []string) error {
				return gsfx.Run(kafka.HelperModule(options))
			},
		}
	}
}
