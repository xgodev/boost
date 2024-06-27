package main

import (
	"fmt"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/log/contrib/rs/zerolog/v1"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {

	zerolog.NewLogger(zerolog.WithLevel("TRACE"))

	basePath, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	log.Infof("current path is %s", basePath)

	collector, err := annotation.Collect(
		annotation.WithFilters("Rest"),
		annotation.WithPath(basePath+"/examples/func/app"),
		annotation.WithPackages("github.com/xgodev/boost/annotation"),
	)
	if err != nil {
		log.Error(err.Error())
	}

	j, _ := yaml.Marshal(collector.Entries())
	fmt.Println(string(j))

}
