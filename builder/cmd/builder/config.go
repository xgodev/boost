package main

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root   = "boost.builder"
	path   = root + ".path"
	output = root + ".output"
)

func init() {
	config.Add(path, ".", "directory scanned for @REST, @GRPC, @Function and the inject annotations")
	config.Add(output, "gen", "directory that receives the inject modules and, under builder/, one package per app; must be inside the Go module")
}
