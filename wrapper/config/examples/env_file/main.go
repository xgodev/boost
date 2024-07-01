package main

import (
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
	"log"
	"os"
)

type AppConfig struct {
	Application struct {
		Name       string
		MyNameTest string
	}
}

func init() {

	os.Setenv("APP_APPLICATION_NAME", "app_test_env")
	os.Setenv("APP_APPLICATION_MY-NAME-TEST", "my_name_test_env")
	os.Setenv("CONF", "./wrapper/config/examples/env_file/config.yaml")

	config.Add("app.application.name", "app_test", "name of application")
	config.Add("app.application.myName", "my_name_test", "name of application")
}

func main() {

	boost.Start()

	c := AppConfig{}

	config.UnmarshalWithPath("app", &c)

	log.Printf(c.Application.Name)
	log.Printf(c.Application.MyNameTest)
}
