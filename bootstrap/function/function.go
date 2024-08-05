package function

import (
	"context"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
	"os"

	co "github.com/spf13/cobra"
)

// CmdFunc defines a function that return a command.
type CmdFunc[T any] func(fn Handler[T]) *co.Command

type Function[T any] struct {
	middlewares []middleware.AnyErrorMiddleware[T]
}

func New[T any](m ...middleware.AnyErrorMiddleware[T]) *Function[T] {
	return &Function[T]{middlewares: m}
}

func (f *Function[T]) Run(ctx context.Context, fn Handler[T], c ...CmdFunc[T]) error {

	// TODO: github.com/alecthomas/kong

	var cmds []*co.Command

	wrp := middleware.NewAnyErrorWrapper[T](ctx, "bootstrap", f.middlewares...)
	for _, v := range c {
		cmds = append(cmds, v(Wrapper[T](wrp, fn)))
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
