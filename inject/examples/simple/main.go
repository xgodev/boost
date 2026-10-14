// Command simple wires an fx application from inject annotations.
//
// Regenerate the modules in gen/ with:
//
//	go generate ./inject/examples/simple
package main

//go:generate go run github.com/xgodev/boost/inject/cmd/inject --boost.inject.output=gen --boost.inject.graphviz=simple.gv

import (
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/inject/examples/simple/app"
	"github.com/xgodev/boost/inject/examples/simple/gen"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

func main() {
	boost.Start()

	var board *app.Board
	if err := fx.New(gen.Module(), fx.Populate(&board), fx.NopLogger).Err(); err != nil {
		log.Fatalf("%v", err)
	}

	for _, line := range board.Lines {
		log.Info(line)
	}
}
