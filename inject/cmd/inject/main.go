package main

import (
	"context"
	"fmt"
	"github.com/xgodev/boost/wrapper/log/contrib/rs/zerolog/v1"
	"golang.org/x/tools/go/packages"
	"os"
	"os/exec"

	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	zerolog.NewLogger(zerolog.WithLevel("INFO"))

	basePath, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()

	log.Infof("current path is %s", basePath)

	moduleName, err := getModuleName(basePath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	entries, err := inject.CollectEntries(basePath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	graph, err := inject.NewGraphFromEntries(ctx, entries)
	if err != nil {
		log.Fatalf(err.Error())
	}

	generator := inject.NewGenerator(moduleName, graph)
	err = generator.Generate(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}

	cmd := exec.Command("go", "mod", "tidy")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("go mod tidy failed: %v", err)
	}

	cmd = exec.Command("go", "mod", "vendor")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("go mod vendor failed: %v", err)
	}

}

func getModuleName(basePath string) (string, error) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedModule, Dir: basePath}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		return "", err
	}
	if len(pkgs) == 0 {
		return "", fmt.Errorf("no packages found in %s", basePath)
	}
	if pkgs[0].Module == nil {
		return "", fmt.Errorf("no module information found in %s", basePath)
	}
	return pkgs[0].Module.Path, nil
}
