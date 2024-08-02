package function

import (
	"context"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
	"os"

	co "github.com/spf13/cobra"
)

// CmdFunc defines a function that return a command.
type CmdFunc func(fn interface{}) *co.Command

type Function struct {
	middlewares []middleware.AnyErrorMiddleware[any]
}

func New(m ...middleware.AnyErrorMiddleware[any]) *Function {
	return &Function{middlewares: m}
}

func (f *Function) Run(ctx context.Context, fn Handler, c ...CmdFunc) error {

	// TODO: github.com/alecthomas/kong

	var cmds []*co.Command

	wrp := middleware.NewAnyErrorWrapper[any](ctx, "bootstrap", f.middlewares...)
	for _, v := range c {
		cmds = append(cmds, v(Wrapper(wrp, fn)))
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
