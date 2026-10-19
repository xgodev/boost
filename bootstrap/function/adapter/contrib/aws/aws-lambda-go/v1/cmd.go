package lambda

import (
	awslambda "github.com/aws/aws-lambda-go/lambda"
	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
)

// New returns the CmdFunc that runs the handler as an AWS Lambda function,
// with the options read from config.
func New[T any]() function.CmdFunc[T] {
	return newCmd[T](NewOptions)
}

// NewWithOptions returns the CmdFunc that runs the handler as an AWS Lambda
// function with the given options.
func NewWithOptions[T any](o *Options) function.CmdFunc[T] {
	return newCmd[T](func() (*Options, error) { return o, nil })
}

func newCmd[T any](options func() (*Options, error)) function.CmdFunc[T] {
	return func(fn function.Handler[T]) *co.Command {
		return &co.Command{
			Use:   "lambda",
			Short: "lambda",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				o, err := options()
				if err != nil {
					return err
				}
				// Start serves invocations until the Lambda runtime stops the
				// process; it does not return.
				awslambda.StartWithOptions(NewHandler[T](fn, o).Handle, awslambda.WithContext(cmd.Context()))
				return nil
			},
		}
	}
}
