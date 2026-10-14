// Command inject generates fx modules from @Provide, @Inject and @Invoke
// annotations. Run it from the application with go:generate:
//
//	//go:generate go run github.com/xgodev/boost/inject/cmd/inject --boost.inject.output=gen/inject
package main

import (
	"context"
	"path/filepath"

	"github.com/xgodev/boost"
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
	logger := log.FromContext(ctx)

	src, err := filepath.Abs(config.String(path))
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error resolving %s", config.String(path)))
	}

	out, err := filepath.Abs(config.String(output))
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error resolving %s", config.String(output)))
	}

	modulePath, moduleDir, err := inject.FindModule(src)
	if err != nil {
		return err
	}

	logger.Infof("scanning %s (module %s)", src, modulePath)

	entries, err := inject.CollectEntries(src)
	if err != nil {
		return err
	}

	g, err := inject.NewGraphFromEntries(ctx, entries)
	if err != nil {
		return err
	}

	if gv := config.String(graphviz); gv != "" {
		if err := inject.ExportInjectGraphToGraphviz(g, gv); err != nil {
			return err
		}
	}

	generator := inject.NewGenerator(g, inject.GeneratorOptions{
		ModulePath: modulePath,
		ModuleDir:  moduleDir,
		OutputDir:  out,
	})
	if err := generator.Generate(ctx); err != nil {
		return err
	}

	logger.Infof("generated %d modules in %s", len(g.Vertices), out)
	return nil
}
