package main

import (
	"fmt"
	"os"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/wrapper/log"
	"gopkg.in/yaml.v3"
)

func main() {

	boost.Start()

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	collector, err := annotation.Collect(
		annotation.WithPath(basePath+"/examples/simple/app"),
		annotation.WithPackages("github.com/xgodev/boost/annotation"),
	)
	if err != nil {
		log.Fatal(err)
	}

	j1, _ := yaml.Marshal(collector.Entries())
	fmt.Println(string(j1))

	entries1 := collector.EntriesWith("MyMethodAnnotation")
	j2, _ := yaml.Marshal(entries1)
	fmt.Println(string(j2))

}
