package main

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root     = "boost.inject"
	path     = root + ".path"
	output   = root + ".output"
	graphviz = root + ".graphviz"
)

func init() {
	config.Add(path, ".", "directory scanned for @Provide, @Inject and @Invoke annotations")
	config.Add(output, "gen/inject", "directory that receives the generated fx modules; must be inside the Go module")
	config.Add(graphviz, "", "optional file that receives the dependency graph in Graphviz DOT format")
}
