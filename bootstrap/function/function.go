package function

import (
	"context"
	"os"
	"strings"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
	"github.com/xgodev/boost/model/errors"

	co "github.com/spf13/cobra"
)

// CmdFunc defines a function that return a command.
type CmdFunc[T any] func(fn Handler[T]) *co.Command

type Function[T any] struct {
	middlewares []middleware.AnyErrorMiddleware[T]
}

func New[T any](m ...middleware.AnyErrorMiddleware[T]) (*Function[T], error) {
	var e T

	switch any(e).(type) {
	case []*event.Event, *event.Event:
		// Tipo válido
	default:
		return nil, errors.New("unsupported handler type")
	}

	return &Function[T]{middlewares: m}, nil
}

// Run runs fn through the adapter command the first argument names. With no
// command named, it runs the one boost.bootstrap.function.default names, or
// the only adapter given: an app wires the one adapter config chose, so it
// runs without arguments. It fails with NotValid when there is no adapter,
// or more than one and none is chosen, instead of printing the usage.
func (f *Function[T]) Run(ctx context.Context, fn Handler[T], c ...CmdFunc[T]) error {

	// TODO: github.com/alecthomas/kong

	if len(c) == 0 {
		return errors.NotValidf("the function has no adapter to run it")
	}

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

	args := os.Args[1:]
	if !namesCommand(args, cmds) {
		def := Default()
		if def == "" && len(cmds) == 1 {
			def = cmds[0].Name()
		}
		if def == "" {
			return errors.NotValidf("no adapter chosen among %s: name one as the first argument or set %s",
				strings.Join(commandNames(cmds), ", "), Root+".default")
		}
		rootCmd.SetArgs(append([]string{def}, args...))
	}

	return cobra.RunContext(ctx, rootCmd, cmds...)
}

// namesCommand reports whether the first argument names one of cmds or a
// command cobra adds (help, completion).
func namesCommand(args []string, cmds []*co.Command) bool {
	if len(args) == 0 {
		return false
	}
	if args[0] == "help" || args[0] == "completion" {
		return true
	}
	for _, cmd := range cmds {
		if cmd.Name() == args[0] || cmd.HasAlias(args[0]) {
			return true
		}
	}
	return false
}

func commandNames(cmds []*co.Command) []string {
	names := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		names = append(names, cmd.Name())
	}
	return names
}
