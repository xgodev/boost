package function

import (
	"github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
	"os"

	co "github.com/spf13/cobra"
)

// CmdFunc defines a function that return a command.
// Valid fn signatures are:
// * func()
// * func() error
// * func(context.Context)
// * func(context.Context) protocol.Result
// * func(event.Event)
// * func(event.Event) protocol.Result
// * func(context.Context, event.Event)
// * func(context.Context, event.Event) protocol.Result
// * func(event.Event) *event.Event
// * func(event.Event) (*event.Event, protocol.Result)
// * func(context.Context, event.Event) *event.Event
// * func(context.Context, event.Event) (*event.Event, protocol.Result)
type CmdFunc func(fn interface{}) *co.Command

// Run executes commands with injected fx modules.
// Valid fn signatures are:
// * func()
// * func() error
// * func(context.Context)
// * func(context.Context) protocol.Result
// * func(event.Event)
// * func(event.Event) protocol.Result
// * func(context.Context, event.Event)
// * func(context.Context, event.Event) protocol.Result
// * func(event.Event) *event.Event
// * func(event.Event) (*event.Event, protocol.Result)
// * func(context.Context, event.Event) *event.Event
// * func(context.Context, event.Event) (*event.Event, protocol.Result)
func Run(fn interface{}, c ...CmdFunc) error {

	var cmds []*co.Command

	for _, v := range c {
		cmds = append(cmds, v(fn))
	}

	rootCmd := &co.Command{
		Use:   "bootstrap",
		Short: "bootstrap",
		Long:  "",
	}

	if def := Default(); def != "" {
		cmd, _, err := rootCmd.Find(os.Args[1:])
		if err == nil && cmd.Use == rootCmd.Use {
			args := append([]string{def}, os.Args[1:]...)
			rootCmd.SetArgs(args)
		}
	}

	return cobra.Run(rootCmd, cmds...)
}
