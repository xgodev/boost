package cmd

import (
	cobra "github.com/xgodev/boost/factory/contrib/spf13/v1"
	"os"

	co "github.com/spf13/cobra"
	"go.uber.org/fx"
)

// CmdFunc defines a function that return a command.
type CmdFunc func(fx.Option) *co.Command

// Run executes commands with injected fx modules.
func Run(options fx.Option, c ...CmdFunc) error {

	var cmds []*co.Command

	for _, v := range c {
		cmds = append(cmds, v(options))
	}

	rootCmd := &co.Command{
		Use:   "faas",
		Short: "faas",
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
