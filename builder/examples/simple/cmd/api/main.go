// Command api serves the REST and gRPC doors of app api.
package main

import (
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/builder/examples/simple/gen/builder/api"
	fxfact "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1"
)

func main() {
	boost.Start()
	fxfact.NewApp(api.Modules()...).Run()
}
