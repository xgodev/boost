// Command samename wires two @Config Options declared in packages that share
// a name: inject tells their types apart by import path.
//
// Regenerate the modules in gen/ with:
//
//	go generate ./inject/examples/samename
package main

//go:generate go run github.com/xgodev/boost/inject/cmd/inject --boost.inject.output=gen

import (
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/inject/examples/samename/app"
	"github.com/xgodev/boost/inject/examples/samename/gen"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

func main() {
	boost.Start()

	var subjects *app.Subjects
	if err := fx.New(gen.Module(), fx.Populate(&subjects), fx.NopLogger).Err(); err != nil {
		log.Fatalf("%v", err)
	}

	log.Infof("price received -> %s, price changed -> %s", subjects.Received, subjects.Changed)
}
