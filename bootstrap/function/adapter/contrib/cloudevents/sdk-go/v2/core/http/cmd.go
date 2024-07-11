package http

import (
	"github.com/cloudevents/sdk-go/v2/client"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
)

// New returns CmdFunc for cloudevents command.
func New(opts ...client.Option) function.CmdFunc {
	return func(fn function.Handler) *co.Command {
		return &co.Command{
			Use:   "cloudevents_http",
			Short: "cloudevents_http",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				return Run(fn, opts...)
			},
		}
	}
}
