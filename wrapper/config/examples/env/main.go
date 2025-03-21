package main

import (
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
	"log"
	"os"
)

type AppConfig struct {
	Application struct {
		Name string
	}
}

func init() {
	config.Add("app.application.name", "app_test", "name of application")
}

func main() {

	boost.Start()

	c1 := AppConfig{}
	config.UnmarshalWithPath("app", &c1)

	log.Printf("%s", c1.Application.Name)

	os.Setenv("APP_APPLICATION_NAME", "app_test_env_2")
	config.Load()

	c2 := AppConfig{}
	config.UnmarshalWithPath("app", &c2)

	log.Printf("%s", c2.Application.Name)

	os.Setenv("APP_APPLICATION_NAME", "app_test_env_3")
	config.Load()

	c3 := AppConfig{}
	config.UnmarshalWithPath("app", &c3)

	log.Printf("%s", c3.Application.Name)

}
