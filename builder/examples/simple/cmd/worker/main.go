// Command worker runs the function of app worker.
package main

import (
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/builder/examples/simple/gen/builder/worker"
	fxfact "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1"
)

func main() {
	boost.Start()
	fxfact.NewApp(worker.Modules()...).Run()
}
