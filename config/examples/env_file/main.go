package main

import (
	"log"
	"os"

	"github.com/xgodev/boost/config"
)

type AppConfig struct {
	Application struct {
		Name       string
		MyNameTest string
	}
}

func init() {
	config.Add("app.application.name", "app_test", "name of application")
	config.Add("app.application.myName", "my_name_test", "name of application")
}

func main() {

	os.Setenv("APP_APPLICATION_NAME", "app_test_env")
	os.Setenv("APP_APPLICATION_MY-NAME-TEST", "my_name_test_env")
	os.Setenv("CONF", "./core/config/examples/env_file/config.yaml")

	config.Load()

	c := AppConfig{}

	config.UnmarshalWithPath("app", &c)

	log.Printf(c.Application.Name)
	log.Printf(c.Application.MyNameTest)
}
