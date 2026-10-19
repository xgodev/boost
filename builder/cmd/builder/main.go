// Command builder generates the inject modules, the door modules and the
// main of every app from the annotations of an application, in one pass. Run it
// from the application with go:generate:
//
//	//go:generate go run github.com/xgodev/boost/builder/cmd/builder --boost.builder.output=gen
//
// Nothing is written when an annotation is not valid.
package main

import (
	"context"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/builder"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {
	boost.Start()

	ctx := context.Background()
	if err := run(ctx); err != nil {
		logger := log.FromContext(ctx)
		logger.Errorf(errors.ErrorStack(err))
		logger.Fatalf(err.Error())
	}
}

func run(ctx context.Context) error {
	return builder.Generate(ctx, config.String(path), config.String(output), config.String(cmd))
}
