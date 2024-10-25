package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/model/errors"
	"golang.org/x/tools/go/packages"
	"os"
	"os/exec"

	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	boost.Start()

	basePath, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	basePath = basePath + "/inject/examples/simple"

	ctx := context.Background()

	log.Infof("current path is %s", basePath)

	moduleName, err := moduleName(basePath)
	if err != nil {
		Err(err)
	}

	entries, err := inject.CollectEntries(basePath)
	if err != nil {
		Err(err)
	}

	graph, err := inject.NewGraphFromEntries(ctx, entries)
	if err != nil {
		Err(err)
	}

	if err := inject.ExportInjectGraphToGraphviz(graph, basePath+"/simple.gv"); err != nil {
		Err(err)
	}

	generator := inject.NewGenerator(moduleName, graph)
	err = generator.Generate(ctx)
	if err != nil {
		Err(err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	err = cmd.Run()
	if err != nil {
		Err(errors.Wrap(err, errors.Internalf("go mod tidy failed")))
	}

	cmd = exec.Command("go", "mod", "vendor")
	err = cmd.Run()
	if err != nil {
		Err(errors.Wrap(err, errors.Internalf("go mod vendor failed")))
	}

}

func Err(err error) {
	log.Errorf(errors.ErrorStack(err))
	log.Fatalf(err.Error())
}

func moduleName(basePath string) (string, error) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedModule, Dir: basePath}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		return "", errors.Wrap(err, errors.Internalf("failed to load packages"))
	}
	if len(pkgs) == 0 {
		return "", errors.Internalf("no packages found in %s", basePath)
	}
	if pkgs[0].Module == nil {
		return "", errors.Internalf("no module found in %s", basePath)
	}
	return pkgs[0].Module.Path, nil
}
