// Command builder generates the inject modules and the door modules of
// every app from the annotations of an application, in one pass. Run it
// from the application with go:generate:
//
//	//go:generate go run github.com/xgodev/boost/builder/cmd/builder --boost.builder.output=gen
//
// Nothing is written when an annotation is not valid.
package main

import (
	"context"
	"path/filepath"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/builder"
	"github.com/xgodev/boost/inject"
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
	src, err := filepath.Abs(config.String(path))
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error resolving %s", config.String(path)))
	}

	out, err := filepath.Abs(config.String(output))
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error resolving %s", config.String(output)))
	}

	return generate(ctx, src, out)
}

// generate writes the inject modules and the doors of the annotations under
// src to out. It writes nothing when an annotation is not valid.
func generate(ctx context.Context, src, out string) error {
	logger := log.FromContext(ctx)

	modulePath, moduleDir, err := inject.FindModule(src)
	if err != nil {
		return err
	}

	logger.Infof("scanning %s (module %s)", src, modulePath)

	entries, err := builder.CollectEntries(src)
	if err != nil {
		return err
	}

	g, err := inject.NewGraphFromEntries(ctx, entries)
	if err != nil {
		return err
	}

	spec, err := builder.NewSpec(entries, g)
	if err != nil {
		return err
	}

	opts := inject.GeneratorOptions{ModulePath: modulePath, ModuleDir: moduleDir, OutputDir: out}
	doors := builder.NewGenerator(spec, opts)

	// render the doors before writing anything, so an error leaves out untouched
	if _, err := doors.Render(ctx); err != nil {
		return err
	}
	if err := inject.NewGenerator(g, opts).Generate(ctx); err != nil {
		return err
	}
	if err := doors.Generate(ctx); err != nil {
		return err
	}

	logger.Infof("generated %d inject modules and %d apps in %s", len(g.Vertices), len(spec.Apps), out)
	return nil
}
