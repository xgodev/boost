package service

import (
	"github.com/dapr/go-sdk/service/common"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
)

// New returns CmdFunc for cloudevents command.
func New(service common.Service) function.CmdFunc {
	return func(fn function.Handler) *co.Command {
		return &co.Command{
			Use:   "dapr",
			Short: "dapr",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				helper := NewHelper(service, fn)
				helper.Start()
				return nil
			},
		}
	}
}
