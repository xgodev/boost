package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/wrapper/log"
	"gopkg.in/yaml.v3"
)

type RestResponse struct {
	Code        int    `attr:"code"`
	Type        string `attr:"type"`
	Description string `attr:"description"`
}

func main() {

	boost.Start()

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	collector, err := annotation.Collect(
		annotation.WithFilters("Rest", "Boost", "Inject"),
		annotation.WithPath(basePath+"/examples/decode/app"),
		annotation.WithPackages("github.com/xgodev/boost/annotation"),
	)
	if err != nil {
		log.Fatal(err)
	}

	entries := collector.Entries()

	j, _ := yaml.Marshal(entries)
	fmt.Println(string(j))

	for _, block := range entries {
		for _, annon := range block.Annotations {
			if annon.Name == "RestResponse" {
				rp := RestResponse{}
				if err := annon.Decode(&rp); err != nil {
					log.Fatal(err)
				}

				j, _ := json.Marshal(rp)
				fmt.Println(string(j))
			}
		}
	}

}
