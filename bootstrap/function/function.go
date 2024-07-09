package function

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
	"github.com/xgodev/boost/middleware"
	"os"

	co "github.com/spf13/cobra"
)

// CmdFunc defines a function that return a command.
type CmdFunc func(fn Handler) *co.Command

type Function struct {
	middlewares []middleware.AnyErrorMiddleware[*event.Event]
}

func New(m ...middleware.AnyErrorMiddleware[*event.Event]) *Function {
	return &Function{middlewares: m}
}

func (f *Function) Run(ctx context.Context, fn Handler, c ...CmdFunc) error {

	wrp := middleware.NewAnyErrorWrapper[*event.Event](ctx, "bootstrap", f.middlewares...)

	var cmds []*co.Command

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
