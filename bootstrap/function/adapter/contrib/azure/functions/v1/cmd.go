package azure

import (
	"context"
	stderrors "errors"
	"fmt"
	"net/http"
	"time"

	co "github.com/spf13/cobra"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// New returns the CmdFunc that runs the handler as an Azure Functions custom
// handler (the azure command), with the options read from config.
func New[T any]() function.CmdFunc[T] {
	return newCmd[T](NewOptions)
}

// NewWithOptions returns the CmdFunc that runs the handler as an Azure
// Functions custom handler with the given options.
func NewWithOptions[T any](o *Options) function.CmdFunc[T] {
	return newCmd[T](func() (*Options, error) { return o, nil })
}

func newCmd[T any](options func() (*Options, error)) function.CmdFunc[T] {
	return func(fn function.Handler[T]) *co.Command {
		return &co.Command{
			Use:   "azure",
			Short: "azure",
			Long:  "",
			RunE: func(cmd *co.Command, args []string) error {
				o, err := options()
				if err != nil {
					return err
				}
				return Serve(cmd.Context(), NewHandler[T](fn, o), o.Port)
			},
		}
	}
}

// Serve answers the functions host on port until ctx ends, then waits for
// the invocations in progress.
func Serve(ctx context.Context, h http.Handler, port int) error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()
	log.FromContext(ctx).Infof("azure custom handler listening on :%d", port)

	select {
	case err := <-errc:
		if stderrors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return errors.Wrap(err, errors.Internalf("azure custom handler on :%d: %s", port, err))
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdown); err != nil {
			return errors.Wrap(err, errors.Internalf("stopping the azure custom handler: %s", err))
		}
		return nil
	}
}
