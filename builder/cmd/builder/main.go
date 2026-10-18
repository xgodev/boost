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

	mains, err := filepath.Abs(config.String(cmd))
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error resolving %s", config.String(cmd)))
	}

	return generate(ctx, src, out, mains)
}

// generate writes the inject modules and the doors of the annotations under
// src to out, and the main of every app to cmdDir. It writes nothing when an
// annotation is not valid or the main of an app was written by hand.
func generate(ctx context.Context, src, out, cmdDir string) error {
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

	spec, err := builder.NewSpec(entries, g, builder.LoadPackages(src))
	if err != nil {
		return err
	}

	opts := builder.Options{
		GeneratorOptions: inject.GeneratorOptions{ModulePath: modulePath, ModuleDir: moduleDir, OutputDir: out},
		CmdDir:           cmdDir,
	}
	doors := builder.NewGenerator(spec, opts)

	// render the doors and mains before writing anything, so an error leaves
	// out and cmdDir untouched
	if _, err := doors.Render(ctx); err != nil {
		return err
	}
	if _, err := doors.RenderCmd(ctx); err != nil {
		return err
	}
	if err := inject.NewGenerator(g, opts.GeneratorOptions).Generate(ctx); err != nil {
		return err
	}
	if err := doors.Generate(ctx); err != nil {
		return err
	}

	logger.Infof("generated %d inject modules and %d apps in %s, mains in %s", len(g.Vertices), len(spec.Apps), out, cmdDir)
	return nil
}
